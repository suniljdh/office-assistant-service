package dbutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/crypto/bcrypt"

	//Support Function
	_ "github.com/denisenkom/go-mssqldb"
)

var (
	server = "LOGIN-PC/MYBOOK"
	//port     = 1433
	user     = "logintech"
	password = "login12ka4"
	database = "office_assistant"

	// server = "67.225.148.206"
	// //port     = 1433
	// user     = "logintechno"
	// password = "sQz0n19$"
	// database = "logintec_ICAI"
)

const (
	secretKey = "6368616e676520746869732070617373776f726420746f206120736563726574"
)

//ConnectDB connect to SQL DB
func ConnectDB() (*sql.DB, error) {
	query := url.Values{}
	query.Add("database", database)

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(user, password),
		//Host:     fmt.Sprintf("%s:%d", server, port),
		Path:     server, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	//log.Println(u.String())
	return sql.Open("sqlserver", u.String())
}

// Encrypt password
func Encrypt(password string) string {
	// Encryption Start
	plainPassword := []byte(password)
	key, _ := hex.DecodeString(secretKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusFailedDependency)
		return err.Error()
	}

	nonce := make([]byte, 12)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		// http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return err.Error()
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusNotAcceptable)
		return err.Error()
	}

	cipherText := aesgcm.Seal(nil, nonce, plainPassword, nil)
	return string(cipherText)
	// Encryption End
}

// Decrypt password
func Decrypt(password string) string {
	if key, err := hex.DecodeString(secretKey); err == nil {
		if cipherText, err := hex.DecodeString(password); err == nil {
			nonce := make([]byte, 12)
			if block, err := aes.NewCipher(key); err == nil {
				if aesgcm, err := cipher.NewGCM(block); err == nil {
					if plainText, err := aesgcm.Open(nil, nonce, cipherText, nil); err == nil {
						return string(plainText)
					}
					return err.Error()

				}
				return err.Error()

			}
			return err.Error()

		}
		return err.Error()
	}
	return password
}

// JSONResult returns json response
func JSONResult(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// HashPassword password hashing
func HashPassword(password string) string {
	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14); err == nil {
		return string(hashedPassword)
	}
	return password
}

// IsValidPassword compare hash & password
func IsValidPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
