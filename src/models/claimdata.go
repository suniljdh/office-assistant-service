package models

import (
	"github.com/dgrijalva/jwt-go"
)

//ClaimData custom claim data for jwt
type ClaimData struct {
	DisplayName string `json:"displayname,omitempty"`
	IsAdmin     bool   `json:"isadmin,omitempty"`
	jwt.StandardClaims
}
