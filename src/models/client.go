package models

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"time"
)

//IClientOperation common operation of client
type IClientOperation interface {
	GetAllClientType(db *sql.DB) (*sql.Rows, error)
	SaveClientData(db *sql.DB) *sql.Row
	FetchClientEntry(db *sql.DB, mode string) *sql.Row
	DeleteClientData(db *sql.DB) *sql.Row
	FetchAllClientData(db *sql.DB) (*sql.Rows, error)
}

//ClientType collects client type like huf, proprietor, individual...
type ClientType struct {
	ID        int    `json:"id"`
	ClientTyp string `json:"type"`
	//ShowFields []byte         `xml:"showfields"`
}

//ClientData saves client info to db
type ClientData struct {
	ClientID     int       `json:"clientid"`
	ClientTypeID int       `json:"clienttypeid"`
	FirstName    string    `json:"firstname"`
	MiddleName   string    `json:"middlename"`
	LastName     string    `json:"lastname"`
	CompanyName  string    `json:"companyname"`
	DateOfBirth  time.Time `json:"dob"`
	PanNo        string    `json:"panno"`
	AadharNo     string    `json:"aadharno"`
	GstNo        string    `json:"gstno"`
	Documents    struct {
		Document []struct {
			DocName        string    `json:"documentname"`
			IsMandatory    bool      `json:"mandatory"`
			Submissiondate time.Time `json:"submissiondate"`
		} `json:"document"`
	} `json:"documents,omitempty"`
}

//GetAllClientType get all client types
func (c *ClientType) GetAllClientType(db *sql.DB) (*sql.Rows, error) {
	tsql := fmt.Sprintf("select id,	isnull(client_type,'') [client_type] from oa.ClientType")

	return db.Query(tsql)

}

//SaveClientData saves client to db
func (c *ClientData) SaveClientData(db *sql.DB) *sql.Row {
	if xmlData, err := xml.Marshal(&c); err == nil {
		return db.QueryRowContext(context.TODO(), "oa.SaveClientData", sql.Named("clientXML", xmlData))
	}
	return nil

	//log.Printf("XMLData: %v", string(xmlData))

}

//FetchClientEntry fetch last client entry
func (c *ClientData) FetchClientEntry(db *sql.DB, mode string) *sql.Row {
	return db.QueryRowContext(context.TODO(), "OA.FetchClientData", sql.Named("clientid", &c.ClientID), sql.Named("dir", mode))
}

//DeleteClientData deleted client by id
func (c *ClientData) DeleteClientData(db *sql.DB) *sql.Row {
	return db.QueryRowContext(context.TODO(), "oa.DeleteClientData", sql.Named("clientid", c.ClientID))
}

//FetchAllClientData get all clients
func (c *ClientData) FetchAllClientData(db *sql.DB) (*sql.Rows, error) {
	return db.QueryContext(context.TODO(), "[oa].[FetchAllClientData]")
}
