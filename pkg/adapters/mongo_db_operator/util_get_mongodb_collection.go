package mongo_db_operator

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoDBCollection(dbURL string) (*mongo.Collection, func()) {
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

	collection := client.Database(parsedURL.DBName()).Collection("vehicles")
	return collection, func() {
		client.Disconnect(context.Background())
	}
}
