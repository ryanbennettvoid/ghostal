package mongo_db_operator

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

const DBPassword = "gho_pass"
const DBName = "gho_db"
const DBPort = "27017"

func createMongoContainer(dbUser string) (string, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:7.0.5",
		ExposedPorts: []string{DBPort + "/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_DATABASE":      DBName,
			"MONGO_INITDB_ROOT_USERNAME": dbUser,
			"MONGO_INITDB_ROOT_PASSWORD": DBPassword,
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

type MockCar struct {
	Make  string `bson:"make"`
	Model string `bson:"model"`
	Year  int    `bson:"year"`
	Color string `bson:"color"`
}

func getCollection(dbURL string) (*mongo.Collection, func()) {
	parsedURL, err := ParseMongoURL(dbURL)
	if err != nil {
		panic(err)
	}
	clonedURL := parsedURL.Clone()
	clonedURL.Path = "admin"
	dbURL = clonedURL.String()

	clientOptions := options.Client().ApplyURI(dbURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	collection := client.Database(DBName).Collection("vehicles")
	return collection, func() {
		client.Disconnect(context.Background())
	}
}

func writeSeedData(dbURL string) {

	vehicles := []interface{}{
		MockCar{"Toyota", "Camry", 2022, "Black"},
		MockCar{"Ford", "Mustang", 2021, "Red"},
		MockCar{"Honda", "Civic", 2020, "Blue"},
		MockCar{"Tesla", "Model 3", 2023, "White"},
		MockCar{"Chevrolet", "Impala", 2019, "Silver"},
	}

	collection, cleanup := getCollection(dbURL)
	defer cleanup()

	_, err := collection.InsertMany(context.TODO(), vehicles)
	if err != nil {
		log.Fatal(err)
	}
}

func getNumVehicles(dbURL string) int {
	collection, cleanup := getCollection(dbURL)
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
	dbURL, cleanup := createMongoContainer(dbUser)
	defer cleanup()

	operator, err := CreateMongoDBOperator(dbURL)
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
		collection, cleanup := getCollection(dbURL)
		defer cleanup()
		_, err := collection.DeleteMany(context.Background(), bson.D{{"year", bson.D{{"$lt", 2022}}}})
		assert.NoError(t, err)
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
