package controllers

import (
	model "models"
	"net/http"
	dbutil "utilities"
)

// FetchAllRolesHandler fetch all roles
func FetchAllRolesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbutil.ConnectDB()
	defer db.Close()
	fatal(err)

	var mr model.Role

	result := make([]model.Role, 0)

	rows, err := mr.FetchAllRoles(db)
	fatal(err)

	for rows.Next() {
		err = rows.Scan(&mr.ID, &mr.RoleName)
		fatal(err)
		result = append(result, mr)
	}

	dbutil.JSONResult(result, w)
}
