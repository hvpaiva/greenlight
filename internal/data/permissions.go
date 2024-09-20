package data

import (
	"context"
	"database/sql"
	"slices"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

func (p Permissions) Contains(permission string) bool {
	return slices.Contains(p, permission)
}

type PermissionModel struct {
	DB *sql.DB
}

func (m PermissionModel) GetAllForUser(userId int64) (Permissions, error) {
	query := `
		SELECT p.code
		FROM permissions p
		INNER JOIN users_permissions up ON up.permission_id = p.id
		INNER JOIN users u ON up.user_id = u.id
		WHERE u.id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var permissions Permissions

	for rows.Next() {
		var permission string
		if err = rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (m PermissionModel) AddForUser(userId int64, permissions ...string) error {
	query := `
		INSERT INTO users_permissions
		SELECT $1, p.id
		FROM permissions p
		WHERE p.code = ANY($2)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userId, pq.Array(permissions))
	return err
}
