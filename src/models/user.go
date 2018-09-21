package models

import (
	"context"
	"database/sql"
	"fmt"
)

// IUser new user method
type IUser interface {
	SaveUser(db *sql.DB) *sql.Row
}

//UserCredentials to validate user
type UserCredentials struct {
	UserName    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	IsAdmin     bool   `json:"is_admin"`
}

// User new user
type User struct {
	UserName string `json:"username"`
	LoginID  string `json:"loginid"`
	Password string `json:"password"`
	IsActive bool   `json:"isactive"`
	RoleID   int    `json:"roleid"`
}

//Token token string
type Token struct {
	Token string `json:"access_token"`
}

//Authorize checks user info
func Authorize(db *sql.DB, user UserCredentials) *sql.Row {
	tsql := fmt.Sprintf("select usr,pwd,display_name, is_admin from oa.User_M where usr = '%s' and pwd = '%s' and is_active=%d", user.UserName, user.Password, 1)

	return db.QueryRow(tsql)
}

// SaveUser save new user
func (u *User) SaveUser(db *sql.DB) *sql.Row {
	return db.QueryRowContext(context.TODO(), "oa.SaveNewUser", sql.Named("username", u.UserName), sql.Named("loginid", u.LoginID), sql.Named("password", u.Password), sql.Named("isactive", u.IsActive), sql.Named("roleid", u.RoleID))
}
