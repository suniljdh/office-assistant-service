package models

import (
	"context"
	"database/sql"
	"utilities"
)

// IUser new user method
type IUser interface {
	SaveUser(db *sql.DB) *sql.Row
}

// UserInfo new user
type UserInfo struct {
	LoginID     string `json:"loginid"`
	Password    string `json:"password"`
	DisplayName string `json:"displayname,omitempty"`
	IsAdmin     bool   `json:"isadmin,omitempty"`
	IsActive    bool   `json:"isactive,omitempty"`
	RoleID      int    `json:"roleid,omitempty"`
}

//Token token string
type Token struct {
	Token string `json:"access_token"`
}

//Authorize checks user info
func (u *UserInfo) Authorize(db *sql.DB) *sql.Row {
	return db.QueryRowContext(context.TODO(), "oa.ValidateCredentials", sql.Named("loginid", u.LoginID))
}

// SaveUser save new user
func (u *UserInfo) SaveUser(db *sql.DB) *sql.Row {
	return db.QueryRowContext(context.TODO(), "oa.SaveNewUser", sql.Named("displayname", u.DisplayName), sql.Named("loginid", u.LoginID), sql.Named("password", dbutil.HashPassword(u.Password)), sql.Named("isactive", u.IsActive), sql.Named("roleid", u.RoleID))
}
