package mongo_db_operator

import (
	"context"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	"ghostel/pkg/values"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

func restoreDB(db *mongo.Client, originalDBName, snapshotDBName string) error {
	// NOTE: MongoDB doesn't support renaming databases (?)
	//		 so cloning is used instead

	// backup original
	backupDBName := "temp_emergency_backup_" + originalDBName
	if err := cloneDB(db, originalDBName, backupDBName); err != nil {
		return err
	}
	// drop original
	if err := dropDB(db, originalDBName); err != nil {
		return err
	}
	// copy snapshot to original
	if err := cloneDB(db, snapshotDBName, originalDBName); err != nil {
		return err
	}
	// drop snapshot
	if err := dropDB(db, snapshotDBName); err != nil {
		return err
	}
	// drop backup
	if err := dropDB(db, backupDBName); err != nil {
		return err
	}
	return nil
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
		return nil, err
	}
	list := make(definitions.List, 0)
	for _, d := range databases.Databases {
		if !strings.HasPrefix(d.Name, values.SnapshotDBPrefix) {
			continue
		}
		snapshotDBNameParts, err := utils.ParseSnapshotDBName(d.Name)
		if err != nil {
			return nil, err
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
		return err
	}

	// Iterate over each collection in the source database
	for _, collection := range collections {
		srcColl := srcDB.Collection(collection)
		dstColl := dstDB.Collection(collection)

		// Find all documents in the source collection
		cur, err := srcColl.Find(context.TODO(), bson.D{})
		if err != nil {
			return err
		}

		// Insert each document into the destination collection
		var docs []interface{}
		for cur.Next(context.TODO()) {
			var elem bson.D
			err := cur.Decode(&elem)
			if err != nil {
				return err
			}
			docs = append(docs, elem)
		}

		if err := cur.Err(); err != nil {
			return err
		}

		cur.Close(context.TODO())

		// Perform the insert to the destination collection
		if len(docs) > 0 {
			_, err = dstColl.InsertMany(context.TODO(), docs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
