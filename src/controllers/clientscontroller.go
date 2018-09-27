package controllers

import (
	"encoding/json"
	"fmt"
	model "models"
	"net/http"
	"strconv"
	dbutil "utilities"
)

// SaveClientHandler save client
func SaveClientHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbutil.ConnectDB()
	defer db.Close()
	fatal(err)
	var c model.ClientData

	err = json.NewDecoder(r.Body).Decode(&c)
	fatal(err)

	//log.Printf("Request Data: %v", clientdata)

	row := c.SaveClientData(db)

	if row != nil {
		var resp struct {
			Context    string `json:"context"`
			ReturnCode string `json:"returncode"`
			ReturnMsg  string `json:"returnmsg"`
		}
		row.Scan(&resp.Context, &resp.ReturnCode, &resp.ReturnMsg)
		dbutil.JSONResult(resp, w)
	}
}

// ClientTypeHandler type of client
func ClientTypeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbutil.ConnectDB()
	defer db.Close()
	fatal(err)

	var c model.ClientType

	result := make([]model.ClientType, 0)

	rows, err := c.GetAllClientType(db)
	fatal(err)

	for rows.Next() {
		err = rows.Scan(&c.ID, &c.ClientTyp)
		fatal(err)
		result = append(result, c)
	}

	dbutil.JSONResult(result, w)

}

// FetchClientHandler fetch client by id
func FetchClientHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbutil.ConnectDB()
	defer db.Close()
	fatal(err)

	clientid, err := strconv.Atoi(r.URL.Query()["id"][0])

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	mode := r.URL.Query()["mode"][0]
	fmt.Println(clientid, mode)
	c := model.ClientData{
		ClientID: clientid,
	}

	row := c.FetchClientEntry(db, mode)

	if row != nil {
		err := row.Scan(
			&c.ClientID,
			&c.ClientTypeID,
			&c.FirstName,
			&c.MiddleName,
			&c.LastName,
			&c.CompanyName,
			&c.DateOfBirth,
			&c.PanNo,
			&c.AadharNo,
			&c.GstNo,
		)

		if err != nil {
			// log.Printf("Error: %v", err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		dbutil.JSONResult(c, w)
	}
}

// DeleteClientHandler delete client by id
func DeleteClientHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbutil.ConnectDB()
	defer db.Close()
	fatal(err)

	clientid, err := strconv.Atoi(r.URL.Query()["id"][0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	c := model.ClientData{
		ClientID: clientid,
	}

	row := c.DeleteClientData(db)
	if row != nil {
		var resp struct {
			ID         int    `json:"clientid"`
			ReturnCode string `json:"returncode"`
			ReturnMsg  string `json:"returnmsg"`
		}
		err := row.Scan(&resp.ID, &resp.ReturnCode, &resp.ReturnMsg)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}

		dbutil.JSONResult(resp, w)
	}
}

// FetchAllClientsHandler fetches all clients
func FetchAllClientsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbutil.ConnectDB()
	defer db.Close()
	fatal(err)

	var c model.ClientData

	result := make([]model.ClientData, 0)

	rows, err := c.FetchAllClientData(db)
	fatal(err)

	for rows.Next() {
		err := rows.Scan(
			&c.ClientID,
			&c.ClientTypeID,
			&c.FirstName,
			&c.MiddleName,
			&c.LastName,
			&c.CompanyName,
			&c.DateOfBirth,
			&c.PanNo,
			&c.AadharNo,
			&c.GstNo,
		)
		fatal(err)
		result = append(result, c)
	}

	dbutil.JSONResult(result, w)
}
