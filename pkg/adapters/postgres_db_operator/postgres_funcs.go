package postgres_db_operator

import (
	"database/sql"
	"fmt"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	"ghostel/pkg/values"
	"github.com/lib/pq"
	"strings"
	"time"
)

func terminateConnections(db *sql.DB, targetDB string) error {
	query := "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = $1;"
	_, err := db.Exec(query, targetDB)
	if err != nil {
		return fmt.Errorf("error terminating connections to database: %w", err)
	}
	return nil
}

func renameDB(db *sql.DB, currentName, newName string) error {
	if err := terminateConnections(db, currentName); err != nil {
		return fmt.Errorf("failed to terminate connection: %w", err)
	}
	query := fmt.Sprintf("ALTER DATABASE %s RENAME TO %s", pq.QuoteIdentifier(currentName), pq.QuoteIdentifier(newName))
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error renaming database from %s to %s: %w", currentName, newName, err)
	}
	return nil
}
func dropDB(db *sql.DB, targetDB string) error {
	if err := terminateConnections(db, targetDB); err != nil {
		return fmt.Errorf("failed to terminate connection: %w", err)
	}
	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", pq.QuoteIdentifier(targetDB))
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func listDBs(db *sql.DB) (definitions.List, error) {
	query := "SELECT datname FROM pg_database"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	list := make([]definitions.ListResult, 0)

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if !strings.HasPrefix(dbName, values.SnapshotDBPrefix) {
			continue
		}

		snapshotDBNameParts, err := utils.ParseSnapshotDBName(dbName)
		if err != nil {
			return nil, err
		}

		list = append(list, definitions.ListResult{
			Name:      snapshotDBNameParts.Name,
			DBName:    dbName,
			CreatedAt: snapshotDBNameParts.Timestamp,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return list, nil
}

func restoreDB(db *sql.DB, originalDBName, snapshotDBName string) error {

	// backup original via rename
	backupSnapshotName := "temp_emergency_backup_" + originalDBName
	if err := renameDB(db, originalDBName, backupSnapshotName); err != nil {
		return fmt.Errorf("failed to rename snapshot: %w", err)
	}

	// restore snapshot to original via rename
	if err := renameDB(db, snapshotDBName, originalDBName); err != nil {
		return fmt.Errorf("failed to rename snapshot: %w", err)
	}

	// delete backup
	if err := dropDB(db, backupSnapshotName); err != nil {
		return fmt.Errorf("failed to remove snapshot: %w", err)
	}

	return nil
}

func snapshotDB(db *sql.DB, originalDBName, originalDBOwner, snapshotName string) error {
	// assumes current db connection is not for default DB, not original DB
	if err := terminateConnections(db, originalDBName); err != nil {
		return err
	}
	fullSnapshotName := utils.BuildSnapshotDBName(snapshotName, time.Now())
	query := fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s OWNER %s;", fullSnapshotName, originalDBName, originalDBOwner)
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create snapshot (%s): %w", query, err)
	}
	return nil
}
