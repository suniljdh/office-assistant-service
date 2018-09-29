package controllers

import (
	"encoding/json"
	"fmt"
	model "models"
	"net/http"
	"utilities"
)

// FetchAllPermissionsHandler fetch all the entities for creating role
func FetchAllPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbutil.ConnectDB()
	defer db.Close()
	fatal(err)

	var p model.MenuPermissions

	rows, err := p.FetchAllPermissions(db)
	result := make([]model.MenuPermissions, 0)

	for rows.Next() {
		err := rows.Scan(
			&p.ID,
			&p.Menu,
			&p.Display,
			&p.IsRoot,
			&p.Parent,
			&p.View,
			&p.Add,
			&p.Edit,
			&p.Delete,
			&p.Search,
			&p.Print,
			&p.Mail,
			&p.Settings,
		)
		fatal(err)
		result = append(result, p)
	}
	// permissions := make([]model.PermissionsData, 0)
	// if len(result) > 0 {
	// 	// log.Printf("row: %#+v\n", result)
	// 	rIdx := 0
	// 	for _, root := range result {
	// 		// log.Printf("row: %d \t %#+v\n", i, row)

	// 		if root.IsRoot == true {
	// 			rIdx++
	// 			var perm model.PermissionsData
	// 			perm.RootMenu.Display = root.Display
	// 			for _, child := range result {
	// 				if child.Parent == root.ID {
	// 					if rIdx == 1 {
	// 						child.FirstSibling = rIdx
	// 						rIdx = 0
	// 					}
	// 					perm.Children = append(perm.Children, child)
	// 				}
	// 			}
	// 			permissions = append(permissions, perm)
	// 		}
	// 	}
	// }
	dbutil.JSONResult(result, w)
}

// SavePermissionHandler save role
func SavePermissionHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbutil.ConnectDB()
	defer db.Close()
	fatal(err)
	var p model.PermissionsData
	// log.Printf("Request Json: %v\n", r.Body)
	err = json.NewDecoder(r.Body).Decode(&p.Permissions)
	defer r.Body.Close()
	fatal(err)

	//log.Printf("Request Data: %v", clientdata)

	row := p.SavePermission(db)

	if row != nil {
		// log.Printf("ROW : %#+v", row)
		var resp struct {
			Context    string `json:"context"`
			ReturnCode string `json:"returncode"`
			ReturnMsg  string `json:"returnmsg"`
		}
		if err := row.Scan(&resp.Context, &resp.ReturnCode, &resp.ReturnMsg); err != nil {
			// http.Error(w, err.Error(), http.StatusNoContent)
			fmt.Println("Scan Error: ", err)
			return
		}
		// log.Printf("%+v", resp)
		dbutil.JSONResult(resp, w)
	}
}

// SaveUserHandler save user
func SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbutil.ConnectDB()
	defer db.Close()
	fatal(err)

	var mu model.UserInfo
	// log.Printf("Request Json: %v\n", r.Body)
	err = json.NewDecoder(r.Body).Decode(&mu)
	defer r.Body.Close()
	fatal(err)

	row := mu.SaveUser(db)

	if row != nil {
		// log.Printf("ROW : %#+v", row)
		var resp struct {
			Context    string `json:"context"`
			ReturnCode string `json:"returncode"`
			ReturnMsg  string `json:"returnmsg"`
		}
		if err := row.Scan(&resp.Context, &resp.ReturnCode, &resp.ReturnMsg); err != nil {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}
		// log.Printf("%+v", resp)
		dbutil.JSONResult(resp, w)
	}
}

// FetchUserPermissionsHandler user & permissions
func FetchUserPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbutil.ConnectDB()
	defer db.Close()
	fatal(err)
	loginid := r.URL.Query()["loginid"][0]
	var p model.UserPermission
	// log.Printf("LoginId : %#v\n", loginid)
	rows, err := p.FetchUserPermissions(db, string(loginid))
	result := make([]model.UserPermission, 0)

	for rows.Next() {
		err := rows.Scan(
			&p.EntityPath,
			&p.EntityName,
			&p.DisplayName,
			&p.Entity,
			&p.BitValue,
			&p.View,
			&p.Add,
			&p.Edit,
			&p.Delete,
			&p.Search,
			&p.Print,
			&p.Mail,
			&p.Settings,
		)
		fatal(err)
		result = append(result, p)
	}
	dbutil.JSONResult(result, w)
}

// FetchEntityPermissionsHandler entity & permissions
func FetchEntityPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbutil.ConnectDB()
	defer db.Close()
	fatal(err)

	var a model.ActionsConfig
	json.NewDecoder(r.Body).Decode(&a)
	// log.Printf("LoginId : %#v\n", loginid)

	row := a.FetchActionPermission(db)

	err = row.Scan(
		&a.LoginID,
		&a.EntityRef,
		&a.View,
		&a.Add,
		&a.Edit,
		&a.Delete,
		&a.Search,
		&a.Print,
		&a.Mail,
		&a.Settings,
	)
	fatal(err)

	dbutil.JSONResult(a, w)
}
