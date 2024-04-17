package postgres_db_operator

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

const DBPassword = "gho_pass"
const DBName = "gho_db"
const DBPort = "5432"

func createPostgresContainer(dbName, dbUser, dbPassword string) (string, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15.1-alpine",
		ExposedPorts: []string{DBPort + "/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPassword,
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

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, DBPassword, host, mappedPort.Port(), DBName)
	return dbURL, func() {
		_ = container.Terminate(ctx)
	}
}

func getNumVehicles(dbURL string) int {
	results := PostgresRunQuery(dbURL, `
		SELECT * FROM vehicles
	`)
	return len(results)
}

func TestIntegration_PostgresDBOperator_Lifecycle(t *testing.T) {
	dbUsers := []string{"postgres", "gho_user"}
	for _, dbUser := range dbUsers {
		runTest(t, dbUser)
	}
}

func runTest(t *testing.T, dbUser string) {
	dbURL, cleanup := createPostgresContainer(DBName, dbUser, DBPassword)
	defer cleanup()

	operator, err := CreatePostgresDBOperator(dbURL)
	assert.NoError(t, err)

	WritePostgresSeedData(dbURL, "vehicles")

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
		PostgresRunQuery(dbURL, `
			DELETE FROM vehicles WHERE year < 2022
		`)
	}

	assert.Equal(t, 2, getNumVehicles(dbURL))

	{
		err := operator.Restore("v1", false)
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
		assert.Len(t, allDatabases, 1)
		assert.Equal(t, "v1", allDatabases[0].Name)
	}
}
