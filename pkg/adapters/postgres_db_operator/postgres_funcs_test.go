package postgres_db_operator

import (
	"database/sql"
	"errors"
	"fmt"
	"ghostal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func listAllDatabases(db *sql.DB) ([]string, error) {
	query := "SELECT datname FROM pg_database"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	list := make([]string, 0)

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		list = append(list, dbName)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return list, nil
}

func TestIntegration_PostgresBackup(t *testing.T) {
	dbURL, cleanupContainer := createPostgresContainer("gho_db", "gho_user", "gho_pass")
	defer cleanupContainer()

	parsedURL, err := ParsePostgresURL(dbURL)
	assert.NoError(t, err)

	postgresClient, cleanupConnection, err := createPostgresConnection(parsedURL, true)
	assert.NoError(t, err)
	defer cleanupConnection()

	// write data to DB
	WritePostgresSeedData(dbURL, "vehicles")

	// attempt destructive operation with backup
	didAttemptDrop := false
	err = backupDB(postgresClient, parsedURL.DBName(), func() error {
		if err := dropDB(postgresClient, parsedURL.DBName()); err != nil {
			panic(err)
		}
		didAttemptDrop = true
		return errors.New("test err")
	})
	assert.Error(t, err)
	assert.Equal(t, "test err", err.Error())
	assert.True(t, didAttemptDrop)

	// verify that original DB is intact
	databases, err := listAllDatabases(postgresClient)
	assert.NoError(t, err)
	_, err = utils.Find(databases, func(dbName string) bool {
		return dbName == parsedURL.DBName()
	})
	assert.NoError(t, err)

	items := PostgresRunQuery(dbURL, `
		SELECT * FROM vehicles
	`)
	assert.NoError(t, err)
	assert.EqualValues(t, 5, len(items))
}
