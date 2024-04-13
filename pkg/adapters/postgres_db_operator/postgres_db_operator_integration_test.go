package postgres_db_operator

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

const DBUser = "postgres"
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

func TestIntegration_PostgresDBOperator_Lifecycle(t *testing.T) {
	dbURL, cleanup := createPostgresContainer()
	defer cleanup()

	operator, err := CreatePostgresDBOperator(dbURL)
	assert.NoError(t, err)

	{
		err := operator.Snapshot("v1")
		assert.NoError(t, err)
	}

	{
		allDatabases, err := operator.List()
		assert.NoError(t, err)
		assert.Len(t, allDatabases, 1)
	}

	{
		err := operator.Snapshot("v2")
		assert.NoError(t, err)
	}

	{
		allDatabases, err := operator.List()
		assert.NoError(t, err)
		assert.Len(t, allDatabases, 2)
	}

	{
		err := operator.Restore("v1")
		assert.NoError(t, err)
	}

	{
		err := operator.Delete("v2")
		assert.NoError(t, err)
	}
}
