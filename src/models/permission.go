package models

import (
	"context"
	"database/sql"
	"encoding/xml"
)

// IPermission for Methods
type IPermission interface {
	FetchAllPermissions(db *sql.DB) (*sql.Rows, error)
	SavePermission(db *sql.DB) *sql.Row
}

// PermissionsData all possible permissions
type PermissionsData struct {
	Permissions []MenuPermissions `json:"permissions"`
}

// MenuPermissions all menu with permission code
type MenuPermissions struct {
	RoleName string `json:"rolename" xml:"rolename"`
	ID       int    `json:"id" xml:"id"`
	Menu     string `json:"menu" xml:"menu"`
	Display  string `json:"display" xml:"display"`
	IsRoot   bool   `json:"is_root" xml:"is_root"`
	Parent   int    `json:"parent,omitempty" xml:"parent"`
	View     bool   `json:"view" xml:"view"`
	Add      bool   `json:"add" xml:"add"`
	Edit     bool   `json:"edit" xml:"edit"`
	Delete   bool   `json:"delete" xml:"delete"`
	Search   bool   `json:"search" xml:"search"`
	Print    bool   `json:"print" xml:"print"`
	Mail     bool   `json:"mail" xml:"mail"`
	Settings bool   `json:"settings" xml:"settings"`
}

// FetchAllPermissions all possible permissions
func (p *MenuPermissions) FetchAllPermissions(db *sql.DB) (*sql.Rows, error) {
	return db.QueryContext(context.TODO(), "oa.Permission")
}

// SavePermission save user permission
func (p *PermissionsData) SavePermission(db *sql.DB) *sql.Row {
	if xmlData, err := xml.Marshal(&p); err == nil {
		// log.Printf("XML : %s ", xmlData)
		rolename := &p.Permissions[0].RoleName
		return db.QueryRowContext(context.TODO(), "oa.SavePermission", sql.Named("roleName", *rolename), sql.Named("permissionXML", xmlData))
	}
	return nil
}
