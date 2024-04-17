package mongo_db_operator

import (
	"fmt"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	"ghostel/pkg/values"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBOperator struct {
	mongoURL *MongoURL
}

func CreateMongoDBOperator(dbURL string) (*MongoDBOperator, error) {
	mongoURL, err := ParseMongoURL(dbURL)
	if err != nil {
		return nil, err
	}
	if mongoURL.dbURL.Scheme != "mongodb" {
		return nil, values.UnsupportedURLSchemeError
	}
	return &MongoDBOperator{
		mongoURL: mongoURL,
	}, nil
}

func (mo *MongoDBOperator) connect(useDefault bool) (*mongo.Client, func(), error) {
	return createMongoConnection(mo.mongoURL, useDefault)
}

func (mo *MongoDBOperator) checkSnapshotName(snapshotName string) error {
	if !utils.IsValidSnapshotName(snapshotName) {
		return values.NoSpecialCharsErr
	}
	list, err := mo.List()
	if err != nil {
		return err
	}
	for _, item := range list {
		if item.Name == snapshotName {
			return values.SnapshotNameTakenErr
		}
	}
	return nil
}

func (mo *MongoDBOperator) Snapshot(snapshotName string) error {
	if err := mo.checkSnapshotName(snapshotName); err != nil {
		return fmt.Errorf("failed to check snapshot name: %w", err)
	}
	db, close, err := mo.connect(true)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer close()
	sourceDatabase := mo.mongoURL.DBName()
	destinationDatabase := snapshotName

	return snapshotDB(db, sourceDatabase, destinationDatabase)
}

func (mo *MongoDBOperator) Restore(snapshotName string, fast bool) error {
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
			return restoreDB(db, originalDBName, snapshotDBName, fast)
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
