package taskrun

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/common"
	"github.com/bytebase/bytebase/backend/common/log"
	"github.com/bytebase/bytebase/backend/component/config"
	"github.com/bytebase/bytebase/backend/component/dbfactory"
	"github.com/bytebase/bytebase/backend/component/state"
	api "github.com/bytebase/bytebase/backend/legacyapi"
	"github.com/bytebase/bytebase/backend/plugin/db"
	"github.com/bytebase/bytebase/backend/plugin/db/mysql"
	"github.com/bytebase/bytebase/backend/store"
	"github.com/bytebase/bytebase/backend/store/model"
	"github.com/bytebase/bytebase/backend/utils"
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"
	v1pb "github.com/bytebase/bytebase/proto/generated-go/v1"
)

// Executor is the task executor.
type Executor interface {
	// RunOnce will be called periodically by the scheduler until terminated is true.
	//
	// NOTE
	//
	// 1. It's possible that err could be non-nil while terminated is false, which
	// usually indicates a transient error and will make scheduler retry later.
	// 2. If err is non-nil, then the detail field will be ignored since info is provided in the err.
	// driverCtx is used by the database driver so that we can cancel the query
	// while have the ability to cleanup migration history etc.
	RunOnce(ctx context.Context, driverCtx context.Context, task *store.TaskMessage, taskRunUID int) (terminated bool, result *api.TaskRunResultPayload, err error)
}

// RunExecutorOnce wraps a TaskExecutor.RunOnce call with panic recovery.
func RunExecutorOnce(ctx context.Context, driverCtx context.Context, exec Executor, task *store.TaskMessage, taskRunUID int) (terminated bool, result *api.TaskRunResultPayload, err error) {
	defer func() {
		if r := recover(); r != nil {
			panicErr, ok := r.(error)
			if !ok {
				panicErr = errors.Errorf("%v", r)
			}
			slog.Error("TaskExecutor PANIC RECOVER", log.BBError(panicErr), log.BBStack("panic-stack"))
			terminated = true
			result = nil
			err = errors.Errorf("TaskExecutor PANIC RECOVER, err: %v", panicErr)
		}
	}()

	return exec.RunOnce(ctx, driverCtx, task, taskRunUID)
}

func getMigrationInfo(ctx context.Context, stores *store.Store, profile config.Profile, task *store.TaskMessage, migrationType db.MigrationType, statement string, schemaVersion model.Version) (*db.MigrationInfo, error) {
	instance, err := stores.GetInstanceV2(ctx, &store.FindInstanceMessage{UID: &task.InstanceID})
	if err != nil {
		return nil, err
	}
	database, err := stores.GetDatabaseV2(ctx, &store.FindDatabaseMessage{UID: task.DatabaseID})
	if err != nil {
		return nil, err
	}
	if database == nil {
		return nil, errors.Errorf("database not found")
	}
	environment, err := stores.GetEnvironmentV2(ctx, &store.FindEnvironmentMessage{ResourceID: &database.EffectiveEnvironmentID})
	if err != nil {
		return nil, err
	}

	mi := &db.MigrationInfo{
		InstanceID:     &instance.UID,
		DatabaseID:     &database.UID,
		CreatorID:      task.CreatorID,
		ReleaseVersion: profile.Version,
		Type:           migrationType,
		Version:        schemaVersion,
		Description:    task.Name,
		Environment:    environment.ResourceID,
		Database:       database.DatabaseName,
		Namespace:      database.DatabaseName,
		Payload:        &storepb.InstanceChangeHistoryPayload{},
	}

	plans, err := stores.ListPlans(ctx, &store.FindPlanMessage{PipelineID: &task.PipelineID})
	if err != nil {
		return nil, err
	}
	if len(plans) == 1 {
		planTypes := []store.PlanCheckRunType{store.PlanCheckDatabaseStatementSummaryReport}
		status := []store.PlanCheckRunStatus{store.PlanCheckRunStatusDone}
		runs, err := stores.ListPlanCheckRuns(ctx, &store.FindPlanCheckRunMessage{
			PlanUID: &plans[0].UID,
			Type:    &planTypes,
			Status:  &status,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list plan check runs")
		}
		sort.Slice(runs, func(i, j int) bool {
			return runs[i].UID > runs[j].UID
		})
		foundChangedResources := false
		for _, run := range runs {
			if foundChangedResources {
				break
			}
			if run.Config.InstanceUid != int32(task.InstanceID) {
				continue
			}
			if run.Config.DatabaseName != database.DatabaseName {
				continue
			}
			if run.Result == nil {
				continue
			}
			for _, result := range run.Result.Results {
				if result.Status != storepb.PlanCheckRunResult_Result_SUCCESS {
					continue
				}
				if report := result.GetSqlSummaryReport(); report != nil {
					mi.Payload.ChangedResources = report.ChangedResources
					foundChangedResources = true
					break
				}
			}
		}
	}

	issue, err := stores.GetIssueV2(ctx, &store.FindIssueMessage{PipelineID: &task.PipelineID})
	if err != nil {
		slog.Error("failed to find containing issue", log.BBError(err))
	}
	if issue != nil {
		// Concat issue title and task name as the migration description so that user can see
		// more context of the migration.
		mi.Description = fmt.Sprintf("%s - %s", issue.Title, task.Name)
		mi.IssueUID = &issue.UID
	}

	mi.Source = db.UI
	creator, err := stores.GetUserByID(ctx, task.CreatorID)
	if err != nil {
		// If somehow we unable to find the principal, we just emit the error since it's not
		// critical enough to fail the entire operation.
		slog.Error("Failed to fetch creator for composing the migration info",
			slog.Int("task_id", task.ID),
			log.BBError(err),
		)
	} else {
		mi.Creator = creator.Name
		mi.CreatorID = creator.ID
	}

	statement = strings.TrimSpace(statement)
	// Only baseline and SDL migration can have empty sql statement, which indicates empty database.
	if mi.Type != db.Baseline && mi.Type != db.MigrateSDL && statement == "" {
		return nil, errors.Errorf("empty statement")
	}
	return mi, nil
}

func executeMigration(
	ctx context.Context,
	driverCtx context.Context,
	stores *store.Store,
	dbFactory *dbfactory.DBFactory,
	stateCfg *state.State,
	profile config.Profile,
	task *store.TaskMessage,
	taskRunUID int,
	statement string,
	sheetID *int,
	mi *db.MigrationInfo) (string, string, error) {
	instance, err := stores.GetInstanceV2(ctx, &store.FindInstanceMessage{UID: &task.InstanceID})
	if err != nil {
		return "", "", err
	}
	database, err := stores.GetDatabaseV2(ctx, &store.FindDatabaseMessage{UID: task.DatabaseID})
	if err != nil {
		return "", "", err
	}

	driver, err := dbFactory.GetAdminDatabaseDriver(ctx, instance, database, db.ConnectionContext{})
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to get driver connection for instance %q", instance.ResourceID)
	}
	defer driver.Close(ctx)

	statementRecord, _ := common.TruncateString(statement, common.MaxSheetSize)
	slog.Debug("Start migration...",
		slog.String("instance", instance.ResourceID),
		slog.String("database", database.DatabaseName),
		slog.String("source", string(mi.Source)),
		slog.String("type", string(mi.Type)),
		slog.String("statement", statementRecord),
	)

	var migrationID string
	opts := db.ExecuteOptions{}
	if task.Type == api.TaskDatabaseDataUpdate && (instance.Engine == storepb.Engine_MYSQL || instance.Engine == storepb.Engine_MARIADB) {
		opts.BeginFunc = func(ctx context.Context, conn *sql.Conn) error {
			updatedTask, err := setThreadIDAndStartBinlogCoordinate(ctx, conn, task, stores)
			if err != nil {
				return errors.Wrap(err, "failed to update the task payload for MySQL rollback SQL")
			}
			task = updatedTask
			return nil
		}
	}
	if task.Type == api.TaskDatabaseDataUpdate && instance.Engine == storepb.Engine_ORACLE {
		// getSetOracleTransactionIdFunc will update the task payload to set the Oracle transaction id, we need to re-retrieve the task to store to the RollbackGenerate.
		opts.EndTransactionFunc = getSetOracleTransactionIDFunc(ctx, task, stores)
	}

	if profile.ExecuteDetail && stateCfg != nil {
		switch task.Type {
		case api.TaskDatabaseSchemaUpdate, api.TaskDatabaseDataUpdate:
			switch instance.Engine {
			case storepb.Engine_MYSQL, storepb.Engine_TIDB, storepb.Engine_OCEANBASE, storepb.Engine_STARROCKS, storepb.Engine_DORIS, storepb.Engine_POSTGRES, storepb.Engine_REDSHIFT, storepb.Engine_RISINGWAVE, storepb.Engine_ORACLE, storepb.Engine_DM, storepb.Engine_OCEANBASE_ORACLE:
				opts.ChunkedSubmission = true
				opts.UpdateExecutionStatus = func(detail *v1pb.TaskRun_ExecutionDetail) {
					stateCfg.TaskRunExecutionStatuses.Store(taskRunUID,
						state.TaskRunExecutionStatus{
							ExecutionStatus: v1pb.TaskRun_EXECUTING,
							ExecutionDetail: detail,
							UpdateTime:      time.Now(),
						})
				}
			default:
				// do nothing
			}
		}
	}

	migrationID, schema, err := utils.ExecuteMigrationDefault(ctx, driverCtx, stores, stateCfg, taskRunUID, driver, mi, statement, sheetID, opts)
	if err != nil {
		return "", "", err
	}

	// If the migration is a data migration, enable the rollback SQL generation and the type of the driver is Oracle, we need to get the rollback SQL before the transaction is committed.
	if task.Type == api.TaskDatabaseDataUpdate && instance.Engine == storepb.Engine_ORACLE {
		updatedTask, err := stores.GetTaskV2ByID(ctx, task.ID)
		if err != nil {
			return "", "", errors.Wrapf(err, "cannot get task by id %d", task.ID)
		}
		payload := &api.TaskDatabaseDataUpdatePayload{}
		if err := json.Unmarshal([]byte(updatedTask.Payload), payload); err != nil {
			return "", "", errors.Wrap(err, "invalid database data update payload")
		}
		if payload.RollbackEnabled {
			// The runner will periodically scan the map to generate rollback SQL asynchronously.
			stateCfg.RollbackGenerate.Store(task.ID, updatedTask)
		}
	}

	if task.Type == api.TaskDatabaseDataUpdate && (instance.Engine == storepb.Engine_MYSQL || instance.Engine == storepb.Engine_MARIADB) {
		conn, err := driver.GetDB().Conn(ctx)
		if err != nil {
			return "", "", errors.Wrap(err, "failed to create connection")
		}
		defer conn.Close()
		updatedTask, err := setMigrationIDAndEndBinlogCoordinate(ctx, conn, task, stores, migrationID)
		if err != nil {
			return "", "", errors.Wrap(err, "failed to update the task payload for MySQL rollback SQL")
		}

		payload := &api.TaskDatabaseDataUpdatePayload{}
		if err := json.Unmarshal([]byte(task.Payload), payload); err != nil {
			return "", "", errors.Wrap(err, "invalid database data update payload")
		}
		if payload.RollbackEnabled {
			// The runner will periodically scan the map to generate rollback SQL asynchronously.
			stateCfg.RollbackGenerate.Store(task.ID, updatedTask)
		}
	}

	return migrationID, schema, nil
}

func getSetOracleTransactionIDFunc(ctx context.Context, task *store.TaskMessage, store *store.Store) func(tx *sql.Tx) error {
	return func(tx *sql.Tx) error {
		payload := &api.TaskDatabaseDataUpdatePayload{}
		if err := json.Unmarshal([]byte(task.Payload), payload); err != nil {
			slog.Error("failed to unmarshal task payload", slog.Int("TaskId", task.ID), log.BBError(err))
			return nil
		}
		// Get oracle current transaction id;
		transactionID, err := tx.QueryContext(ctx, "SELECT RAWTOHEX(tx.xid) FROM v$transaction tx JOIN v$session s ON tx.ses_addr = s.saddr")
		if err != nil {
			slog.Error("failed to transaction id in task", slog.Int("TaskId", task.ID), log.BBError(err))
			return nil
		}
		defer transactionID.Close()
		var txID string
		for transactionID.Next() {
			err := transactionID.Scan(&txID)
			if err != nil {
				slog.Error("failed to the Oracle transaction id in task", slog.Int("TaskId", task.ID), log.BBError(err))
				return nil
			}
		}
		if err := transactionID.Err(); err != nil {
			return err
		}
		payload.TransactionID = txID
		updatedPayload, err := json.Marshal(payload)
		if err != nil {
			slog.Error("failed to unmarshal task payload", slog.Int("TaskId", task.ID), log.BBError(err), slog.Any("payload", updatedPayload))
			return nil
		}
		updatedPayloadString := string(updatedPayload)
		patch := &api.TaskPatch{
			ID:        task.ID,
			UpdaterID: api.SystemBotID,
			Payload:   &updatedPayloadString,
		}
		if _, err = store.UpdateTaskV2(ctx, patch); err != nil {
			slog.Error("failed to update task with new payload", slog.Any("TaskPatch", patch), log.BBError(err))
			return nil
		}
		return nil
	}
}

func setThreadIDAndStartBinlogCoordinate(ctx context.Context, conn *sql.Conn, task *store.TaskMessage, store *store.Store) (*store.TaskMessage, error) {
	payload := &api.TaskDatabaseDataUpdatePayload{}
	if err := json.Unmarshal([]byte(task.Payload), payload); err != nil {
		return nil, errors.Wrap(err, "invalid database data update payload")
	}

	var connID string
	if err := conn.QueryRowContext(ctx, "SELECT CONNECTION_ID();").Scan(&connID); err != nil {
		return nil, errors.Wrap(err, "failed to get the connection ID")
	}
	payload.ThreadID = connID

	binlogInfo, err := mysql.GetBinlogInfo(ctx, conn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the binlog info before executing the migration transaction")
	}
	if (binlogInfo == api.BinlogInfo{}) {
		slog.Warn("binlog is not enabled", slog.Int("task", task.ID))
		return task, nil
	}
	payload.BinlogFileStart = binlogInfo.FileName
	payload.BinlogPosStart = binlogInfo.Position

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal task payload")
	}
	payloadString := string(payloadBytes)
	patch := &api.TaskPatch{
		ID:        task.ID,
		UpdaterID: api.SystemBotID,
		Payload:   &payloadString,
	}
	updatedTask, err := store.UpdateTaskV2(ctx, patch)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to patch task %d with the MySQL thread ID", task.ID)
	}
	return updatedTask, nil
}

func setMigrationIDAndEndBinlogCoordinate(ctx context.Context, conn *sql.Conn, task *store.TaskMessage, store *store.Store, migrationID string) (*store.TaskMessage, error) {
	payload := &api.TaskDatabaseDataUpdatePayload{}
	if err := json.Unmarshal([]byte(task.Payload), payload); err != nil {
		return nil, errors.Wrap(err, "invalid database data update payload")
	}

	payload.MigrationID = migrationID
	binlogInfo, err := mysql.GetBinlogInfo(ctx, conn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the binlog info before executing the migration transaction")
	}
	if (binlogInfo == api.BinlogInfo{}) {
		slog.Warn("binlog is not enabled", slog.Int("task", task.ID))
		return task, nil
	}
	payload.BinlogFileEnd = binlogInfo.FileName
	payload.BinlogPosEnd = binlogInfo.Position

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal task payload")
	}
	payloadString := string(payloadBytes)
	patch := &api.TaskPatch{
		ID:        task.ID,
		UpdaterID: api.SystemBotID,
		Payload:   &payloadString,
	}
	updatedTask, err := store.UpdateTaskV2(ctx, patch)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to patch task %d with the MySQL thread ID", task.ID)
	}
	return updatedTask, nil
}

func postMigration(ctx context.Context, stores *store.Store, task *store.TaskMessage, mi *db.MigrationInfo, migrationID string, sheetID *int) (bool, *api.TaskRunResultPayload, error) {
	instance, err := stores.GetInstanceV2(ctx, &store.FindInstanceMessage{UID: &task.InstanceID})
	if err != nil {
		return true, nil, err
	}
	database, err := stores.GetDatabaseV2(ctx, &store.FindDatabaseMessage{UID: task.DatabaseID})
	if err != nil {
		return true, nil, err
	}

	if mi.Type == db.Migrate || mi.Type == db.MigrateSDL {
		if _, err := stores.UpdateDatabase(ctx, &store.UpdateDatabaseMessage{
			InstanceID:    instance.ResourceID,
			DatabaseName:  database.DatabaseName,
			SchemaVersion: &mi.Version,
		}, api.SystemBotID); err != nil {
			return true, nil, errors.Errorf("failed to update database %q for instance %q", database.DatabaseName, instance.ResourceID)
		}
	}

	slog.Debug("Post migration...",
		slog.String("instance", instance.ResourceID),
		slog.String("database", database.DatabaseName),
	)

	// Set schema config.
	if sheetID != nil && task.DatabaseID != nil {
		sheet, err := stores.GetSheet(ctx, &store.FindSheetMessage{
			UID: sheetID,
		})
		if err != nil {
			slog.Error("Failed to get sheet from store", slog.Int("sheetID", *sheetID), log.BBError(err))
		} else if sheet.Payload != nil && (sheet.Payload.DatabaseConfig != nil || sheet.Payload.BaselineDatabaseConfig != nil) {
			databaseSchema, err := stores.GetDBSchema(ctx, *task.DatabaseID)
			if err != nil {
				slog.Error("Failed to get database config from store", slog.Int("sheetID", *sheetID), slog.Int("databaseUID", *task.DatabaseID), log.BBError(err))
			} else {
				updatedDatabaseConfig := utils.MergeDatabaseConfig(sheet.Payload.BaselineDatabaseConfig, databaseSchema.GetConfig(), sheet.Payload.DatabaseConfig)
				err = stores.UpdateDBSchema(ctx, *task.DatabaseID, &store.UpdateDBSchemaMessage{
					Config: updatedDatabaseConfig,
				}, api.SystemBotID)
				if err != nil {
					slog.Error("Failed to update database config", slog.Int("sheetID", *sheetID), slog.Int("databaseUID", *task.DatabaseID), log.BBError(err))
				}
			}
		}
	}

	// Remove schema drift anomalies.
	if err := stores.ArchiveAnomalyV2(ctx, &store.ArchiveAnomalyMessage{
		DatabaseUID: task.DatabaseID,
		Type:        api.AnomalyDatabaseSchemaDrift,
	}); err != nil && common.ErrorCode(err) != common.NotFound {
		slog.Error("Failed to archive anomaly",
			slog.String("instance", instance.ResourceID),
			slog.String("database", database.DatabaseName),
			slog.String("type", string(api.AnomalyDatabaseSchemaDrift)),
			log.BBError(err))
	}

	detail := fmt.Sprintf("Applied migration version %s to database %q.", mi.Version.Version, database.DatabaseName)
	if mi.Type == db.Baseline {
		detail = fmt.Sprintf("Established baseline version %s for database %q.", mi.Version.Version, database.DatabaseName)
	}

	storedVersion, err := mi.Version.Marshal()
	if err != nil {
		slog.Error("failed to convert database schema version",
			slog.String("version", mi.Version.Version),
			log.BBError(err),
		)
	}
	return true, &api.TaskRunResultPayload{
		Detail:        detail,
		MigrationID:   migrationID,
		ChangeHistory: fmt.Sprintf("instances/%s/databases/%s/changeHistories/%s", instance.ResourceID, database.DatabaseName, migrationID),
		Version:       storedVersion,
	}, nil
}

func runMigration(ctx context.Context, driverCtx context.Context, store *store.Store, dbFactory *dbfactory.DBFactory, stateCfg *state.State, profile config.Profile, task *store.TaskMessage, taskRunUID int, migrationType db.MigrationType, statement string, schemaVersion model.Version, sheetID *int) (terminated bool, result *api.TaskRunResultPayload, err error) {
	mi, err := getMigrationInfo(ctx, store, profile, task, migrationType, statement, schemaVersion)
	if err != nil {
		return true, nil, err
	}

	migrationID, _, err := executeMigration(ctx, driverCtx, store, dbFactory, stateCfg, profile, task, taskRunUID, statement, sheetID, mi)
	if err != nil {
		return true, nil, err
	}
	return postMigration(ctx, store, task, mi, migrationID, sheetID)
}
