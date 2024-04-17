package mongo_db_operator

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

const DBPassword = "gho_pass"
const DBName = "gho_db"
const DBPort = "27017"

func createMongoContainer(dbName, dbUser, dbPassword string) (string, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:7.0.5",
		ExposedPorts: []string{DBPort + "/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_DATABASE":      dbName,
			"MONGO_INITDB_ROOT_USERNAME": dbUser,
			"MONGO_INITDB_ROOT_PASSWORD": dbPassword,
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

	dbURL := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?tls=false", dbUser, DBPassword, host, mappedPort.Port(), DBName)
	return dbURL, func() {
		_ = container.Terminate(ctx)
	}
}

func getNumVehicles(dbURL string) int {
	collection, cleanup := GetMongoDBCollection(dbURL, "vehicles")
	defer cleanup()
	result, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}
	return int(result)
}

func TestIntegration_MongoDBOperator_Lifecycle(t *testing.T) {
	dbUsers := []string{"admin", "gho_user"}
	for _, dbUser := range dbUsers {
		runTest(t, dbUser)
	}
}

func runTest(t *testing.T, dbUser string) {
	dbURL, cleanup := createMongoContainer(DBName, dbUser, DBPassword)
	defer cleanup()

	operator, err := CreateMongoDBOperator(dbURL)
	assert.NoError(t, err)

	WriteMongoDBSeedData(dbURL, "vehicles")

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
		collection, cleanup := GetMongoDBCollection(dbURL, "vehicles")
		defer cleanup()
		_, err := collection.DeleteMany(context.Background(), bson.D{{"year", bson.D{{"$lt", 2022}}}})
		assert.NoError(t, err)
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
