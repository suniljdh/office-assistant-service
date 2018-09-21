package controllers

import (
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	model "models"
	"net/http"
	"strings"
	"time"
	"utilities"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

const (
	privKeyPath = "keys/app.rsa"
	pubKeyPath  = "keys/app.rsa.pub"
)

var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
)

var (
	db *sql.DB
	// tmpl *template.Template
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// INITKey initialization of private and public keys
func INITKey() {
	var err error
	signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)

}

// LoginHandler userlogin
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// tmpl.Execute(w, nil)

	db, err := dbutil.ConnectDB()
	fatal(err)
	var u model.UserCredentials

	err = json.NewDecoder(r.Body).Decode(&u)
	// u.UserName = r.FormValue("username")
	// u.Password = r.FormValue("password")
	// log.Printf("UserName: %s\nPassword: %s", &u.UserName, &u.Password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in Request")
		return
	}
	//defer db.Close()
	// err = db.Ping()
	// if err != nil {
	// 	log.Printf("From Ping() Attempt: %s" + err.Error())
	// }

	fatal(err)
	row := model.Authorize(db, u)

	var userData model.UserCredentials

	row.Scan(&userData.UserName, &userData.Password, &userData.DisplayName, &userData.IsAdmin)

	// log.Printf("%v", userData)

	if strings.ToLower(u.UserName) != userData.UserName || u.Password != userData.Password {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		fmt.Fprint(w, "Invalid credentials")
		return
	}
	claims := model.ClaimData{
		userData.DisplayName,
		userData.IsAdmin,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}

	tokenData := model.Token{Token: tokenString}
	dbutil.JSONResult(&tokenData, w)
}

// ValidateTokenMiddleware token valid or not
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
}
