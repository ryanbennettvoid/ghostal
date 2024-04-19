package mongo_db_operator

import (
	"fmt"
	"ghostal/pkg/definitions"
	"ghostal/pkg/utils"
	"ghostal/pkg/values"
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
	list, err := mo.ListSnapshots()
	if err != nil {
		return err
	}
	_, err = utils.Find(list, func(item definitions.SnapshotListResult) bool {
		return item.SnapshotName == snapshotName
	})
	if err == nil {
		// item found
		return values.SnapshotNameTakenErr
	}
	return nil
}

func (mo *MongoDBOperator) Snapshot(snapshotName string) error {
	if err := mo.checkSnapshotName(snapshotName); err != nil {
		return err
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

	allDatabases, err := listSnapshots(db, mo.mongoURL.DBName())
	if err != nil {
		return err
	}
	for _, d := range allDatabases {
		if d.SnapshotName == snapshotName {
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

	allDatabases, err := listSnapshots(db, mo.mongoURL.DBName())
	if err != nil {
		return err
	}
	for _, d := range allDatabases {
		if d.SnapshotName == snapshotName {
			snapshotDBName := d.DBName
			return dropDB(db, snapshotDBName)
		}
	}

	return values.SnapshotNotExistsErr
}

func (mo *MongoDBOperator) ListSnapshots() (definitions.SnapshotList, error) {
	db, close, err := mo.connect(true)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer close()

	return listSnapshots(db, mo.mongoURL.DBName())
}
