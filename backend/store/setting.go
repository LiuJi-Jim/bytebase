package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/common"
	api "github.com/bytebase/bytebase/backend/legacyapi"
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"
)

// FindSettingMessage is the message for finding setting.
type FindSettingMessage struct {
	Name    *api.SettingName
	Enforce bool
}

// SetSettingMessage is the message for updating setting.
type SetSettingMessage struct {
	Name        api.SettingName
	Value       string
	Description *string
}

// SettingMessage is the message of setting.
type SettingMessage struct {
	Name        api.SettingName
	Value       string
	Description string
	CreatedTs   int64
}

// GetWorkspaceGeneralSetting gets the workspace general setting payload.
func (s *Store) GetWorkspaceGeneralSetting(ctx context.Context) (*storepb.WorkspaceProfileSetting, error) {
	settingName := api.SettingWorkspaceProfile
	setting, err := s.GetSettingV2(ctx, &FindSettingMessage{
		Name:    &settingName,
		Enforce: true,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get setting %s", settingName)
	}
	if setting == nil {
		return nil, errors.Errorf("cannot find setting %v", settingName)
	}

	payload := new(storepb.WorkspaceProfileSetting)
	if err := common.ProtojsonUnmarshaler.Unmarshal([]byte(setting.Value), payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func (s *Store) GetAppIMSetting(ctx context.Context) (*storepb.AppIMSetting, error) {
	settingName := api.SettingAppIM
	setting, err := s.GetSettingV2(ctx, &FindSettingMessage{
		Name: &settingName,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get setting %s", settingName)
	}
	if setting == nil {
		return nil, errors.Errorf("cannot find setting %v", settingName)
	}

	payload := new(storepb.AppIMSetting)
	if err := common.ProtojsonUnmarshaler.Unmarshal([]byte(setting.Value), payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// GetWorkspaceID finds the workspace id in setting bb.workspace.id.
func (s *Store) GetWorkspaceID(ctx context.Context) (string, error) {
	settingName := api.SettingWorkspaceID
	setting, err := s.GetSettingV2(ctx, &FindSettingMessage{
		Name: &settingName,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to get setting %s", settingName)
	}
	if setting == nil {
		return "", errors.Errorf("cannot find setting %v", settingName)
	}
	return setting.Value, nil
}

// GetWorkspaceApprovalSetting gets the workspace approval setting.
func (s *Store) GetWorkspaceApprovalSetting(ctx context.Context) (*storepb.WorkspaceApprovalSetting, error) {
	settingName := api.SettingWorkspaceApproval
	setting, err := s.GetSettingV2(ctx, &FindSettingMessage{
		Name: &settingName,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get setting %s", settingName)
	}
	if setting == nil {
		return nil, errors.Errorf("cannot find setting %v", settingName)
	}

	payload := new(storepb.WorkspaceApprovalSetting)
	if err := common.ProtojsonUnmarshaler.Unmarshal([]byte(setting.Value), payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// GetWorkspaceExternalApprovalSetting gets the workspace external approval setting.
func (s *Store) GetWorkspaceExternalApprovalSetting(ctx context.Context) (*storepb.ExternalApprovalSetting, error) {
	settingName := api.SettingWorkspaceExternalApproval
	setting, err := s.GetSettingV2(ctx, &FindSettingMessage{
		Name: &settingName,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get setting %s", settingName)
	}
	if setting == nil {
		return nil, errors.Errorf("cannot find setting %v", settingName)
	}

	payload := new(storepb.ExternalApprovalSetting)
	if err := common.ProtojsonUnmarshaler.Unmarshal([]byte(setting.Value), payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// GetMaskingAlgorithmSetting gets the masking algorithm setting.
func (s *Store) GetMaskingAlgorithmSetting(ctx context.Context) (*storepb.MaskingAlgorithmSetting, error) {
	settingName := api.SettingMaskingAlgorithm
	setting, err := s.GetSettingV2(ctx, &FindSettingMessage{
		Name: &settingName,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get setting %s", settingName)
	}
	if setting == nil {
		return &storepb.MaskingAlgorithmSetting{}, nil
	}

	payload := new(storepb.MaskingAlgorithmSetting)
	if err := common.ProtojsonUnmarshaler.Unmarshal([]byte(setting.Value), payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// GetSemanticTypesSetting gets the semantic types setting.
func (s *Store) GetSemanticTypesSetting(ctx context.Context) (*storepb.SemanticTypeSetting, error) {
	settingName := api.SettingSemanticTypes
	setting, err := s.GetSettingV2(ctx, &FindSettingMessage{
		Name: &settingName,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get setting %s", settingName)
	}
	if setting == nil {
		return &storepb.SemanticTypeSetting{}, nil
	}

	payload := new(storepb.SemanticTypeSetting)
	if err := common.ProtojsonUnmarshaler.Unmarshal([]byte(setting.Value), payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// GetDataClassificationSetting gets the data classification setting.
func (s *Store) GetDataClassificationSetting(ctx context.Context) (*storepb.DataClassificationSetting, error) {
	settingName := api.SettingDataClassification
	setting, err := s.GetSettingV2(ctx, &FindSettingMessage{
		Name: &settingName,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get setting %s", settingName)
	}
	if setting == nil {
		return &storepb.DataClassificationSetting{}, nil
	}

	payload := new(storepb.DataClassificationSetting)
	if err := common.ProtojsonUnmarshaler.Unmarshal([]byte(setting.Value), payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// GetDataClassificationConfigByID gets the classification config by the id.
func (s *Store) GetDataClassificationConfigByID(ctx context.Context, classificationConfigID string) (*storepb.DataClassificationSetting_DataClassificationConfig, error) {
	setting, err := s.GetDataClassificationSetting(ctx)
	if err != nil {
		return nil, err
	}
	for _, config := range setting.Configs {
		if config.Id == classificationConfigID {
			return config, nil
		}
	}
	return &storepb.DataClassificationSetting_DataClassificationConfig{}, nil
}

// DeleteCache deletes the cache.
func (s *Store) DeleteCache() {
	s.settingCache.Purge()
	s.policyCache.Purge()
	s.userEmailCache.Purge()
	s.userIDCache.Purge()
}

// GetSettingV2 returns the setting by name.
func (s *Store) GetSettingV2(ctx context.Context, find *FindSettingMessage) (*SettingMessage, error) {
	if find.Name != nil && !find.Enforce {
		if v, ok := s.settingCache.Get(*find.Name); ok {
			return v, nil
		}
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	settings, err := listSettingV2Impl(ctx, tx, find)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list setting")
	}
	if len(settings) == 0 {
		return nil, nil
	}
	if len(settings) > 1 {
		return nil, errors.Errorf("found multiple settings: %v", find.Name)
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}
	return settings[0], nil
}

// ListSettingV2 returns a list of settings.
func (s *Store) ListSettingV2(ctx context.Context, find *FindSettingMessage) ([]*SettingMessage, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()
	settings, err := listSettingV2Impl(ctx, tx, find)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list setting")
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}

	for _, setting := range settings {
		s.settingCache.Add(setting.Name, setting)
	}
	return settings, nil
}

// UpsertSettingV2 upserts the setting by name.
func (s *Store) UpsertSettingV2(ctx context.Context, update *SetSettingMessage, principalUID int) (*SettingMessage, error) {
	fields := []string{"creator_id", "updater_id", "updated_ts", "name", "value"}
	updateFields := []string{"value = EXCLUDED.value", "updater_id = EXCLUDED.updater_id", "updated_ts = EXCLUDED.updated_ts"}
	valuePlaceholders, args := []string{"$1", "$2", "$3", "$4", "$5"}, []any{principalUID, principalUID, time.Now().Unix(), update.Name, update.Value}

	if v := update.Description; v != nil {
		fields = append(fields, "description")
		valuePlaceholders = append(valuePlaceholders, fmt.Sprintf("$%d", len(args)+1))
		updateFields = append(updateFields, "description = EXCLUDED.description")
		args = append(args, *v)
	}
	query := `INSERT INTO setting (` + strings.Join(fields, ", ") + `) 
		VALUES (` + strings.Join(valuePlaceholders, ", ") + `) 
		ON CONFLICT (name) DO UPDATE SET ` + strings.Join(updateFields, ", ") + `
		RETURNING name, value, description, created_ts`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	var setting SettingMessage
	if err := tx.QueryRowContext(ctx, query, args...).Scan(
		&setting.Name,
		&setting.Value,
		&setting.Description,
		&setting.CreatedTs,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, &common.Error{Code: common.NotFound, Err: errors.Errorf("setting not found: %s", update.Name)}
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}
	s.settingCache.Add(setting.Name, &setting)
	return &setting, nil
}

// CreateSettingIfNotExistV2 creates a new setting only if the named setting doesn't exist.
func (s *Store) CreateSettingIfNotExistV2(ctx context.Context, create *SettingMessage, principalUID int) (*SettingMessage, bool, error) {
	if v, ok := s.settingCache.Get(create.Name); ok {
		return v, false, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()
	settings, err := listSettingV2Impl(ctx, tx, &FindSettingMessage{Name: &create.Name})
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to list settings")
	}
	if len(settings) > 1 {
		return nil, false, errors.Errorf("found settings for setting name: %v", create.Name)
	}
	if len(settings) == 1 {
		// Don't create setting if the named setting already exists.
		return settings[0], false, nil
	}

	fields := []string{"creator_id", "updater_id", "name", "value", "description"}
	valuesPlaceholders, args := []string{"$1", "$2", "$3", "$4", "$5"}, []any{principalUID, principalUID, create.Name, create.Value, create.Description}

	query := `INSERT INTO setting (` + strings.Join(fields, ",") + `)
		VALUES (` + strings.Join(valuesPlaceholders, ",") + `)
		RETURNING name, value, description, created_ts`
	var setting SettingMessage
	if err := tx.QueryRowContext(ctx, query, args...).Scan(
		&setting.Name,
		&setting.Value,
		&setting.Description,
		&setting.CreatedTs,
	); err != nil {
		return nil, false, err
	}

	if err := tx.Commit(); err != nil {
		return nil, false, errors.Wrap(err, "failed to commit transaction")
	}
	s.settingCache.Add(setting.Name, &setting)
	return &setting, true, nil
}

// DeleteSettingV2 deletes a setting by the name.
func (s *Store) DeleteSettingV2(ctx context.Context, name api.SettingName) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM setting WHERE name = $1`, name); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	s.settingCache.Remove(name)
	return nil
}

func listSettingV2Impl(ctx context.Context, tx *Tx, find *FindSettingMessage) ([]*SettingMessage, error) {
	where, args := []string{"TRUE"}, []any{}
	if v := find.Name; v != nil {
		where, args = append(where, fmt.Sprintf("name = $%d", len(args)+1)), append(args, *v)
	}
	rows, err := tx.QueryContext(ctx, `
		SELECT
			name,
			value,
			description,
			created_ts
		FROM setting
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settingMessages []*SettingMessage
	for rows.Next() {
		var settingMessage SettingMessage
		if err := rows.Scan(
			&settingMessage.Name,
			&settingMessage.Value,
			&settingMessage.Description,
			&settingMessage.CreatedTs,
		); err != nil {
			return nil, err
		}
		settingMessages = append(settingMessages, &settingMessage)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return settingMessages, nil
}
