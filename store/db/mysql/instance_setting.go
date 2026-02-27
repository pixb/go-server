package mysql

import (
	"context"

	"github.com/pixb/go-server/store"
)

func (d *Driver) UpsertInstanceSetting(ctx context.Context, upsert *store.InstanceSetting) (*store.InstanceSetting, error) {
	stmt := `
		INSERT INTO system_setting (
			name, value, description
		)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			value = VALUES(value),
			description = VALUES(description)
	`
	if _, err := d.db.ExecContext(ctx, stmt, upsert.Name, upsert.Value, upsert.Description); err != nil {
		return nil, err
	}

	return upsert, nil
}

func (d *Driver) ListInstanceSettings(ctx context.Context, find *store.FindInstanceSetting) ([]*store.InstanceSetting, error) {
	query := `
		SELECT
			name,
			value,
			description
		FROM system_setting
		WHERE 1 = 1`
	args := []interface{}{}

	if find.Name != "" {
		query += " AND name = ?"
		args = append(args, find.Name)
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*store.InstanceSetting{}
	for rows.Next() {
		systemSettingMessage := &store.InstanceSetting{}
		if err := rows.Scan(
			&systemSettingMessage.Name,
			&systemSettingMessage.Value,
			&systemSettingMessage.Description,
		); err != nil {
			return nil, err
		}
		list = append(list, systemSettingMessage)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *Driver) DeleteInstanceSetting(ctx context.Context, delete *store.DeleteInstanceSetting) error {
	stmt := "DELETE FROM system_setting WHERE name = ?"
	_, err := d.db.ExecContext(ctx, stmt, delete.Name)
	return err
}
