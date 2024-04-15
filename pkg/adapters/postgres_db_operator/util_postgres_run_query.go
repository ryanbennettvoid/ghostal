package postgres_db_operator

import (
	"database/sql"
	"fmt"
)

// PostgresRunQuery executes a query and returns the results as a slice of maps where each map represents a row.
func PostgresRunQuery(dbURL, query string) []map[string]interface{} {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(fmt.Errorf("failed to open database: %v", err))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %v", err))
	}

	rows, err := db.Query(query)
	if err != nil {
		panic(fmt.Errorf("failed to execute query: %v", err))
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		panic(fmt.Errorf("failed to get columns: %v", err))
	}

	// Prepare a slice to hold the values and a slice of interfaces for scanning.
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	var results []map[string]interface{}

	// Fetch rows and scan into a map.
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			panic(fmt.Errorf("failed to scan row: %v", err))
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := valuePtrs[i].(*interface{})
			rowMap[col] = *val
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		panic(fmt.Errorf("error processing rows: %v", err))
	}

	return results
}
