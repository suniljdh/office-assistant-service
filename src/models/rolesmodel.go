package models

import (
	"context"
	"database/sql"
)

// IRole all roles related methods
type IRole interface {
	FetchAllRoles(db *sql.DB) (*sql.Rows, error)
}

// Role role structure
type Role struct {
	ID       int    `json:"id"`
	RoleName string `json:"role_name"`
}

// FetchAllRoles all possible roles
func (r *Role) FetchAllRoles(db *sql.DB) (*sql.Rows, error) {
	return db.QueryContext(context.TODO(), "oa.FetchAllRoles")
}
