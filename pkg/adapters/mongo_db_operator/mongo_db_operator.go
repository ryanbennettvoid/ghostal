package mongo_db_operator

import (
	"context"
	"fmt"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	"ghostel/pkg/values"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBOperator struct {
	mongoURL *MongoURL
}

func CreateMongoDBOperator(dbURL string) (*MongoDBOperator, error) {
	mongoURL, err := ParseMongoURL(dbURL)
	if err != nil {
		return nil, err
	}
	return &MongoDBOperator{
		mongoURL: mongoURL,
	}, nil
}

func (mo *MongoDBOperator) SupportsDatabase(dbURL string) (bool, error) {
	scheme, err := utils.GetURLScheme(dbURL)
	if err != nil {
		return false, fmt.Errorf("failed to get URL scheme: %w", err)
	}
	return scheme == "mongodb", nil
}

func (mo *MongoDBOperator) connect(useDefault bool) (*mongo.Client, func(), error) {
	dbURL := mo.mongoURL.dbURL.String()
	if useDefault {
		newMongoURL := mo.mongoURL.Clone()
		newMongoURL.Path = "admin"
		dbURL = newMongoURL.String()
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURL))
	if err != nil {
		return nil, nil, err
	}
	return client, func() {
		_ = client.Disconnect(context.Background())
	}, nil
}

func (mo *MongoDBOperator) Snapshot(snapshotName string) error {
	db, close, err := mo.connect(true)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer close()

	sourceDatabase := mo.mongoURL.DBName()
	destinationDatabase := snapshotName

	return snapshotDB(db, sourceDatabase, destinationDatabase)
}

func (mo *MongoDBOperator) Restore(snapshotName string) error {
	db, close, err := mo.connect(true)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer close()

	allDatabases, err := listDBs(db)
	if err != nil {
		return err
	}
	for _, d := range allDatabases {
		if d.Name == snapshotName {
			originalDBName := mo.mongoURL.DBName()
			snapshotDBName := d.DBName
			return restoreDB(db, originalDBName, snapshotDBName)
		}
	}

	return values.SnapshotNotExistsErr
}

func (mo *MongoDBOperator) Delete(snapshotName string) error {
	db, close, err := mo.connect(true)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer close()

	allDatabases, err := listDBs(db)
	if err != nil {
		return err
	}
	for _, d := range allDatabases {
		if d.Name == snapshotName {
			snapshotDBName := d.DBName
			return dropDB(db, snapshotDBName)
		}
	}

	return values.SnapshotNotExistsErr
}

func (mo *MongoDBOperator) List() (definitions.List, error) {
	db, close, err := mo.connect(true)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer close()

	return listDBs(db)
}
