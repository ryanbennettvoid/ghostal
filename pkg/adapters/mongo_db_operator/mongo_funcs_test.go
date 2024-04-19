package mongo_db_operator

import (
	"context"
	"errors"
	"ghostal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestIntegration_MongoBackup(t *testing.T) {
	dbURL, cleanupContainer := createMongoContainer("gho_db", "gho_user", "gho_pass")
	defer cleanupContainer()

	parsedURL, err := ParseMongoURL(dbURL)
	assert.NoError(t, err)

	mongoClient, cleanupConnection, err := createMongoConnection(parsedURL, true)
	assert.NoError(t, err)
	defer cleanupConnection()

	// write data to DB
	WriteMongoDBSeedData(dbURL, "vehicles")

	// attempt destructive operation with backup
	didAttemptDrop := false
	err = backupDB(mongoClient, parsedURL.DBName(), func() error {
		if err := dropDB(mongoClient, parsedURL.DBName()); err != nil {
			panic(err)
		}
		didAttemptDrop = true
		return errors.New("test err")
	})
	assert.Error(t, err)
	assert.Equal(t, "test err", err.Error())
	assert.True(t, didAttemptDrop)

	// verify that original DB is intact
	databases, err := mongoClient.ListDatabases(context.Background(), bson.D{})
	assert.NoError(t, err)
	_, err = utils.Find(databases.Databases, func(db mongo.DatabaseSpecification) bool {
		return db.Name == parsedURL.DBName()
	})
	assert.NoError(t, err)

	collection := mongoClient.Database(parsedURL.DBName()).Collection("vehicles")
	numItems, err := collection.CountDocuments(context.Background(), bson.D{})
	assert.NoError(t, err)
	assert.EqualValues(t, 5, numItems)
}
