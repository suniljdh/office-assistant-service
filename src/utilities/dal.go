package dbutil

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"

	//Support Function
	_ "github.com/denisenkom/go-mssqldb"
)

var (
	server = "LOGIN-PC/MYBOOK"
	//port     = 1433
	user     = "logintech"
	password = "login12ka4"
	database = "office_assistant"
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
