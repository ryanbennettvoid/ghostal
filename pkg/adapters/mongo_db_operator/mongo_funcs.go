package mongo_db_operator

import (
	"context"
	"errors"
	"fmt"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	"ghostel/pkg/values"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

func createMongoConnection(mongoURL *MongoURL, useDefault bool) (*mongo.Client, func(), error) {
	dbURL := mongoURL.dbURL.String()
	if useDefault {
		newMongoURL := mongoURL.Clone()
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

// backupDB backs up `sourceDB` and restores it if `fn` fails
func backupDB(db *mongo.Client, sourceDB string, fn func() error) error {
	backupDBName := "temp_emergency_backup_" + sourceDB
	if err := cloneDB(db, sourceDB, backupDBName); err != nil {
		return fmt.Errorf("failed clone original to backup: %w", err)
	}
	if err := fn(); err != nil {
		// if error, drop current source and rename backup to source
		_ = dropDB(db, sourceDB)
		_ = cloneDB(db, backupDBName, sourceDB)
		// after emergency restore, drop backup
		_ = dropDB(db, backupDBName)
		return err
	}
	// is success, drop backup
	_ = dropDB(db, backupDBName)
	return nil
}

func restoreDB(db *mongo.Client, originalDBName, snapshotDBName string, fast bool) error {
	// NOTE: MongoDB doesn't support renaming databases (?)
	//		 so cloning is used instead

	if fast {
		// drop original
		if err := dropDB(db, originalDBName); err != nil {
			return fmt.Errorf("failed to drop original: %w", err)
		}
		// copy snapshot to original
		if err := cloneDB(db, snapshotDBName, originalDBName); err != nil {
			return fmt.Errorf("failed to clone snapshot to orignal: %w", err)
		}
		return nil
	}

	return backupDB(db, originalDBName, func() error {
		// drop original
		if err := dropDB(db, originalDBName); err != nil {
			return fmt.Errorf("failed to drop original: %w", err)
		}
		// copy snapshot to original
		if err := cloneDB(db, snapshotDBName, originalDBName); err != nil {
			return fmt.Errorf("failed to clone snapshot to orignal: %w", err)
		}
		return nil
	})
}

func snapshotDB(db *mongo.Client, originalDBName, snapshotName string) error {
	fullSnapshotName := utils.BuildSnapshotDBName(snapshotName, time.Now())
	return cloneDB(db, originalDBName, fullSnapshotName)
}

func dropDB(db *mongo.Client, dbName string) error {
	return db.Database(dbName).Drop(context.Background())
}

func listDBs(db *mongo.Client) (definitions.List, error) {
	// List all collections in the source database
	databases, err := db.ListDatabases(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}
	list := make(definitions.List, 0)
	for _, d := range databases.Databases {
		if !strings.HasPrefix(d.Name, values.SnapshotDBPrefix) {
			continue
		}
		snapshotDBNameParts, err := utils.ParseSnapshotDBName(d.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to parse snapshot database name: %w", err)
		}
		list = append(list, definitions.ListResult{
			Name:      snapshotDBNameParts.Name,
			DBName:    d.Name,
			CreatedAt: snapshotDBNameParts.Timestamp,
		})
	}
	return list, nil
}

func cloneDB(db *mongo.Client, sourceDBName, targetDBName string) error {

	srcDB := db.Database(sourceDBName)
	dstDB := db.Database(targetDBName)

	// List all collections in the source database
	collections, err := srcDB.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		return fmt.Errorf("failed to list collection names: %w", err)
	}

	if len(collections) == 0 {
		return errors.New("cannot clone: source database has no collections")
	}

	// Iterate over each collection in the source database
	for _, collection := range collections {
		srcColl := srcDB.Collection(collection)
		dstColl := dstDB.Collection(collection)

		// Find all documents in the source collection
		cur, err := srcColl.Find(context.TODO(), bson.D{})
		if err != nil {
			return fmt.Errorf("failed to find documents: %w", err)
		}

		// Insert each document into the destination collection
		var docs []interface{}
		for cur.Next(context.TODO()) {
			var elem bson.D
			err := cur.Decode(&elem)
			if err != nil {
				return fmt.Errorf("failed to decode document: %w", err)
			}
			docs = append(docs, elem)
		}

		if err := cur.Err(); err != nil {
			return fmt.Errorf("cursor error: %s", err)
		}

		cur.Close(context.TODO())

		// Perform the insert to the destination collection
		if len(docs) > 0 {
			_, err = dstColl.InsertMany(context.TODO(), docs)
			if err != nil {
				return fmt.Errorf("failed to insert many: %w", err)
			}
		}
	}

	return nil
}
