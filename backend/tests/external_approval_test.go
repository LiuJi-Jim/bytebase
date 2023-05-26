package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	api "github.com/bytebase/bytebase/backend/legacyapi"
	"github.com/bytebase/bytebase/backend/plugin/db"
	"github.com/bytebase/bytebase/backend/tests/fake"
	v1pb "github.com/bytebase/bytebase/proto/generated-go/v1"
)

func TestExternalApprovalFeishu_AllUserCanBeFound(t *testing.T) {
	t.Parallel()
	a := require.New(t)
	ctx := context.Background()
	ctl := &controller{}
	dataDir := t.TempDir()
	ctx, err := ctl.StartServerWithExternalPg(ctx, &config{
		dataDir:                 dataDir,
		vcsProviderCreator:      fake.NewGitLab,
		feishuProverdierCreator: fake.NewFeishu,
	})
	a.NoError(err)
	defer ctl.Close(ctx)

	err = ctl.setLicense()
	a.NoError(err)

	// close existing issues
	issues, err := ctl.getIssues(nil /* projectID */)
	a.NoError(err)
	for _, issue := range issues {
		patchedIssue, err := ctl.patchIssueStatus(api.IssueStatusPatch{
			ID:     issue.ID,
			Status: api.IssueCanceled,
		})
		a.NoError(err)
		a.Equal(api.IssueCanceled, patchedIssue.Status)
	}

	_, err = ctl.settingServiceClient.SetSetting(ctx, &v1pb.SetSettingRequest{
		Setting: &v1pb.Setting{
			Name: fmt.Sprintf("settings/%s", api.SettingAppIM),
			Value: &v1pb.Value{
				Value: &v1pb.Value_AppImSettingValue{
					AppImSettingValue: &v1pb.AppIMSetting{
						ImType:    v1pb.AppIMSetting_FEISHU,
						AppId:     "123",
						AppSecret: "123",
						ExternalApproval: &v1pb.AppIMSetting_ExternalApproval{
							Enabled: true,
						},
					},
				},
			},
		},
	})
	a.NoError(err)

	// Create a DBA account.
	dbaUser, err := ctl.authServiceClient.CreateUser(ctx, &v1pb.CreateUserRequest{
		User: &v1pb.User{
			Title:    "DBA",
			Email:    "dba@dba.com",
			UserRole: v1pb.UserRole_DBA,
			UserType: v1pb.UserType_USER,
			Password: "dbapass",
		},
	})
	a.NoError(err)
	dbaUserUID, err := strconv.Atoi(strings.TrimPrefix(dbaUser.Name, "users/"))
	a.NoError(err)

	err = ctl.feishuProvider.RegisterEmails("demo@example.com", "dba@dba.com")
	a.NoError(err)

	// Create a project.
	project, err := ctl.createProject(ctx)
	a.NoError(err)
	projectUID, err := strconv.Atoi(project.Uid)
	a.NoError(err)

	// Provision an instance.
	instanceRootDir := t.TempDir()
	instanceName := "testInstance1"
	instanceDir, err := ctl.provisionSQLiteInstance(instanceRootDir, instanceName)
	a.NoError(err)

	_, prodEnvironmentUID, err := ctl.getEnvironment(ctx, "prod")
	a.NoError(err)

	// Add an instance.
	instance, err := ctl.addInstance(api.InstanceCreate{
		ResourceID:    generateRandomString("instance", 10),
		EnvironmentID: prodEnvironmentUID,
		Name:          instanceName,
		Engine:        db.SQLite,
		Host:          instanceDir,
	})
	a.NoError(err)

	// Expecting project to have no database.
	databases, err := ctl.getDatabases(api.DatabaseFind{
		ProjectID: &projectUID,
	})
	a.NoError(err)
	a.Zero(len(databases))
	// Expecting instance to have no database.
	databases, err = ctl.getDatabases(api.DatabaseFind{
		InstanceID: &instance.ID,
	})
	a.NoError(err)
	a.Zero(len(databases))

	// Create an issue that creates a database.
	databaseName := "testSchemaUpdate"
	err = ctl.createDatabase(ctx, projectUID, instance, databaseName, "", nil /* labelMap */)
	a.NoError(err)

	// Expecting project to have 1 database.
	databases, err = ctl.getDatabases(api.DatabaseFind{
		ProjectID: &projectUID,
	})
	a.NoError(err)
	a.Equal(1, len(databases))
	database := databases[0]
	a.Equal(instance.ID, database.Instance.ID)

	sheet, err := ctl.createSheet(api.SheetCreate{
		ProjectID:  projectUID,
		Name:       "migration statement sheet",
		Statement:  migrationStatement,
		Visibility: api.ProjectSheet,
		Source:     api.SheetFromBytebaseArtifact,
		Type:       api.SheetForSQL,
	})
	a.NoError(err)

	// Create an issue that updates database schema.
	createContext, err := json.Marshal(&api.MigrationContext{
		DetailList: []*api.MigrationDetail{
			{
				MigrationType: db.Migrate,
				DatabaseID:    database.ID,
				SheetID:       sheet.ID,
			},
		},
	})
	a.NoError(err)
	issue, err := ctl.createIssue(api.IssueCreate{
		ProjectID:     projectUID,
		Name:          fmt.Sprintf("update schema for database %q", databaseName),
		Type:          api.IssueDatabaseSchemaUpdate,
		Description:   fmt.Sprintf("This updates the schema of database %q.", databaseName),
		AssigneeID:    dbaUserUID,
		CreateContext: string(createContext),
	})
	a.NoError(err)

	for {
		review, err := ctl.reviewServiceClient.GetReview(ctx, &v1pb.GetReviewRequest{
			Name: fmt.Sprintf("projects/%d/reviews/%d", issue.ProjectID, issue.ID),
		})
		a.NoError(err)
		if review.ApprovalFindingDone {
			break
		}
		time.Sleep(time.Second)
	}

	attention := true
	issue, err = ctl.patchIssue(issue.ID, api.IssuePatch{
		AssigneeNeedAttention: &attention,
	})
	a.NoError(err)
	a.Equal(true, issue.AssigneeNeedAttention)

	// Sleep for a few seconds, giving time to ApplicationRunner to create external approvals.
	time.Sleep(ctl.profile.AppRunnerInterval + 2*time.Second)
	issue, err = ctl.getIssue(issue.ID)
	a.NoError(err)
	taskStatus, err := getNextTaskStatus(issue)
	a.NoError(err)
	// The task is still waiting for approval.
	a.Equal(api.TaskPendingApproval, taskStatus)

	// Should have 1 PENDING approval on the feishu side.
	a.Equal(1, ctl.feishuProvider.PendingApprovalCount())
	// Simulate users approving on the feishu side.
	ctl.feishuProvider.ApprovePendingApprovals()

	// Waiting ApplicationRunner to approves the issue.
	status, err := ctl.waitIssuePipelineWithNoApproval(ctx, issue.ID)
	a.NoError(err)
	a.Equal(api.TaskDone, status)
}

func TestExternalApprovalFeishu_AssigneeCanBeFound(t *testing.T) {
	t.Parallel()
	a := require.New(t)
	ctx := context.Background()
	ctl := &controller{}
	dataDir := t.TempDir()
	ctx, err := ctl.StartServerWithExternalPg(ctx, &config{
		dataDir:                 dataDir,
		vcsProviderCreator:      fake.NewGitLab,
		feishuProverdierCreator: fake.NewFeishu,
	})
	a.NoError(err)
	defer ctl.Close(ctx)

	err = ctl.setLicense()
	a.NoError(err)

	// close existing issues
	issues, err := ctl.getIssues(nil /* projectID */)
	a.NoError(err)
	for _, issue := range issues {
		patchedIssue, err := ctl.patchIssueStatus(api.IssueStatusPatch{
			ID:     issue.ID,
			Status: api.IssueCanceled,
		})
		a.NoError(err)
		a.Equal(api.IssueCanceled, patchedIssue.Status)
	}

	_, err = ctl.settingServiceClient.SetSetting(ctx, &v1pb.SetSettingRequest{
		Setting: &v1pb.Setting{
			Name: fmt.Sprintf("settings/%s", api.SettingAppIM),
			Value: &v1pb.Value{
				Value: &v1pb.Value_AppImSettingValue{
					AppImSettingValue: &v1pb.AppIMSetting{
						ImType:    v1pb.AppIMSetting_FEISHU,
						AppId:     "123",
						AppSecret: "123",
						ExternalApproval: &v1pb.AppIMSetting_ExternalApproval{
							Enabled: true,
						},
					},
				},
			},
		},
	})
	a.NoError(err)

	// Create a DBA account.
	// Create a DBA account.
	dbaUser, err := ctl.authServiceClient.CreateUser(ctx, &v1pb.CreateUserRequest{
		User: &v1pb.User{
			Title:    "DBA",
			Email:    "dba@dba.com",
			UserRole: v1pb.UserRole_DBA,
			UserType: v1pb.UserType_USER,
			Password: "dbapass",
		},
	})
	a.NoError(err)
	dbaUserUID, err := strconv.Atoi(strings.TrimPrefix(dbaUser.Name, "users/"))
	a.NoError(err)

	err = ctl.feishuProvider.RegisterEmails("dba@dba.com")
	a.NoError(err)

	// Create a project.
	project, err := ctl.createProject(ctx)
	a.NoError(err)
	projectUID, err := strconv.Atoi(project.Uid)
	a.NoError(err)

	// Provision an instance.
	instanceRootDir := t.TempDir()
	instanceName := "testInstance1"
	instanceDir, err := ctl.provisionSQLiteInstance(instanceRootDir, instanceName)
	a.NoError(err)

	_, prodEnvironmentUID, err := ctl.getEnvironment(ctx, "prod")
	a.NoError(err)

	// Add an instance.
	instance, err := ctl.addInstance(api.InstanceCreate{
		ResourceID:    generateRandomString("instance", 10),
		EnvironmentID: prodEnvironmentUID,
		Name:          instanceName,
		Engine:        db.SQLite,
		Host:          instanceDir,
	})
	a.NoError(err)

	// Expecting project to have no database.
	databases, err := ctl.getDatabases(api.DatabaseFind{
		ProjectID: &projectUID,
	})
	a.NoError(err)
	a.Zero(len(databases))
	// Expecting instance to have no database.
	databases, err = ctl.getDatabases(api.DatabaseFind{
		InstanceID: &instance.ID,
	})
	a.NoError(err)
	a.Zero(len(databases))

	// Create an issue that creates a database.
	databaseName := "testSchemaUpdate"
	err = ctl.createDatabase(ctx, projectUID, instance, databaseName, "", nil /* labelMap */)
	a.NoError(err)

	// Expecting project to have 1 database.
	databases, err = ctl.getDatabases(api.DatabaseFind{
		ProjectID: &projectUID,
	})
	a.NoError(err)
	a.Equal(1, len(databases))
	database := databases[0]
	a.Equal(instance.ID, database.Instance.ID)

	sheet, err := ctl.createSheet(api.SheetCreate{
		ProjectID:  projectUID,
		Name:       "migration statement sheet",
		Statement:  migrationStatement,
		Visibility: api.ProjectSheet,
		Source:     api.SheetFromBytebaseArtifact,
		Type:       api.SheetForSQL,
	})
	a.NoError(err)

	// Create an issue that updates database schema.
	createContext, err := json.Marshal(&api.MigrationContext{
		DetailList: []*api.MigrationDetail{
			{
				MigrationType: db.Migrate,
				DatabaseID:    database.ID,
				SheetID:       sheet.ID,
			},
		},
	})
	a.NoError(err)
	issue, err := ctl.createIssue(api.IssueCreate{
		ProjectID:     projectUID,
		Name:          fmt.Sprintf("update schema for database %q", databaseName),
		Type:          api.IssueDatabaseSchemaUpdate,
		Description:   fmt.Sprintf("This updates the schema of database %q.", databaseName),
		AssigneeID:    dbaUserUID,
		CreateContext: string(createContext),
	})
	a.NoError(err)

	for {
		review, err := ctl.reviewServiceClient.GetReview(ctx, &v1pb.GetReviewRequest{
			Name: fmt.Sprintf("projects/%d/reviews/%d", issue.ProjectID, issue.ID),
		})
		a.NoError(err)
		if review.ApprovalFindingDone {
			break
		}
		time.Sleep(time.Second)
	}

	attention := true
	issue, err = ctl.patchIssue(issue.ID, api.IssuePatch{
		AssigneeNeedAttention: &attention,
	})
	a.NoError(err)
	a.Equal(true, issue.AssigneeNeedAttention)

	// Sleep for a few seconds, giving time to ApplicationRunner to create external approvals.
	time.Sleep(ctl.profile.AppRunnerInterval + 2*time.Second)
	issue, err = ctl.getIssue(issue.ID)
	a.NoError(err)
	taskStatus, err := getNextTaskStatus(issue)
	a.NoError(err)
	// The task is still waiting for approval.
	a.Equal(api.TaskPendingApproval, taskStatus)

	// Should have 1 PENDING approval on the feishu side.
	a.Equal(1, ctl.feishuProvider.PendingApprovalCount())
	// Simulate users approving on the feishu side.
	ctl.feishuProvider.ApprovePendingApprovals()

	// Waiting ApplicationRunner to approves the issue.
	status, err := ctl.waitIssuePipelineWithNoApproval(ctx, issue.ID)
	a.NoError(err)
	a.Equal(api.TaskDone, status)
}
