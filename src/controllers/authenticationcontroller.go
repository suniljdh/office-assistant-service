package controllers

import (
	"crypto/rsa"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	model "models"
	"net/http"
	"text/template"
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
	db   *sql.DB
	tmpl *template.Template
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

	db, err := dbutil.ConnectDB()
	fatal(err)

	var u model.UserInfo

	// if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
	// 	http.Error(w, err.Error(), http.StatusUnauthorized)
	// 	return
	// }

	u.LoginID = r.FormValue("txtLoginId")
	u.Password = r.FormValue("txtPassword")

	// log.Printf("User Info : %#+v\n", u)

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("From Ping() Attempt: %s" + err.Error())
		fatal(err)
	}
	plainPassword := u.Password

	row := u.Authorize(db)

	if err := row.Scan(
		&u.LoginID,
		&u.Password,
		&u.DisplayName,
		&u.IsAdmin,
		&u.RoleID,
	); err != nil {
		// log.Printf("Login Error: %#+v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	isvalid := dbutil.IsValidPassword(u.Password, plainPassword)

	if !isvalid {
		http.Error(w, "Unauthorized Access", http.StatusUnauthorized)
		return
	}

	claims := model.ClaimData{
		u.DisplayName,
		u.IsAdmin,
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

	expire := time.Now().Add(time.Minute * time.Duration(5))
	cookie := http.Cookie{
		Name:    "accesstoken",
		Value:   tokenString,
		Expires: expire,
	}

	http.SetCookie(w, &cookie)

	// tokenData := model.Token{Token: tokenString}
	// dbutil.JSONResult(&tokenData, w)
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

// IndexHandler home page handler
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseGlob("views/*.html"))

	tmpl.ExecuteTemplate(w, "index.html", nil)
}
