package postgres_db_operator

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

const DBUser = "gho_user"
const DBPassword = "gho_pass"
const DBName = "gho_db"
const DBPort = "5432"

func createPostgresContainer() (string, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15.1-alpine",
		ExposedPorts: []string{DBPort + "/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       DBName,
			"POSTGRES_USER":     DBUser,
			"POSTGRES_PASSWORD": DBPassword,
		},
		WaitingFor: wait.ForListeningPort(DBPort + "/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to start container: %s", err))
	}

	host, err := container.Host(context.Background())
	if err != nil {
		panic(err)
	}

	mappedPort, err := container.MappedPort(context.Background(), DBPort)
	if err != nil {
		panic(err)
	}

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", DBUser, DBPassword, host, mappedPort.Port(), DBName)
	return dbURL, func() {
		_ = container.Terminate(ctx)
	}
}

// runQuery executes a query and returns the results as a slice of maps where each map represents a row.
func runQuery(dbURL, query string) []map[string]interface{} {
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

func writeSeedData(dbURL string) {
	runQuery(dbURL, `
		CREATE TABLE vehicles (
			id SERIAL PRIMARY KEY,
			make VARCHAR(50),
			model VARCHAR(50),
			year INT,
			color VARCHAR(50)
		);
	`)
	runQuery(dbURL, `
		INSERT INTO vehicles (make, model, year, color) VALUES
		('Toyota', 'Camry', 2022, 'Black'),
		('Ford', 'Mustang', 2021, 'Red'),
		('Honda', 'Civic', 2020, 'Blue'),
		('Tesla', 'Model 3', 2023, 'White'),
		('Chevrolet', 'Impala', 2019, 'Silver');
	`)
}

func getNumVehicles(dbURL string) int {
	results := runQuery(dbURL, `
		SELECT * FROM vehicles
	`)
	return len(results)
}

func TestIntegration_PostgresDBOperator_Lifecycle(t *testing.T) {
	dbURL, cleanup := createPostgresContainer()
	defer cleanup()

	operator, err := CreatePostgresDBOperator(dbURL)
	assert.NoError(t, err)

	writeSeedData(dbURL)

	assert.Equal(t, 5, getNumVehicles(dbURL))

	{
		assert.NoError(t, operator.Snapshot("v1"))
	}

	{
		allDatabases, err := operator.List()
		assert.NoError(t, err)
		assert.Len(t, allDatabases, 1)
	}

	{
		assert.NoError(t, operator.Snapshot("v2"))
	}

	{
		allDatabases, err := operator.List()
		assert.NoError(t, err)
		assert.Len(t, allDatabases, 2)
	}

	assert.Equal(t, 5, getNumVehicles(dbURL))

	{
		// modify DB before restoring snapshot
		runQuery(dbURL, `
			DELETE FROM vehicles WHERE year < 2022
		`)
	}

	assert.Equal(t, 2, getNumVehicles(dbURL))

	{
		err := operator.Restore("v1")
		assert.NoError(t, err)
	}

	assert.Equal(t, 5, getNumVehicles(dbURL))

	{
		err := operator.Delete("v2")
		assert.NoError(t, err)
	}

	{
		allDatabases, err := operator.List()
		assert.NoError(t, err)
		assert.Len(t, allDatabases, 0)
	}
}
