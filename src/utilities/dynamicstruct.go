package dbutil

import (
	"database/sql"
	"log"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// CreateStruct from Sql Rows
func CreateStruct(rows *sql.Rows) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	cols, err := rows.Columns()
	fatal(err)

	columns := make([]interface{}, len(cols))
	columnspointer := make([]interface{}, len(cols))

	for i := range columns {
		columnspointer[i] = &columns[i]
	}
	for rows.Next() {
		rows.Scan(columnspointer...)
		row := make(map[string]interface{})
		for n, colName := range cols {
			val := columnspointer[n].(*interface{})
			row[colName] = *val
		}

		result = append(result, row)
	}

	return result
}
