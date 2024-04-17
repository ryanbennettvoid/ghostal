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

func createPostgresConnection(postgresURL *PostgresURL, useDefault bool) (*sql.DB, func(), error) {
	dbURL := postgresURL.dbURL.String()
	if useDefault {
		newPGURL := postgresURL.Clone()
		newPGURL.Path = "postgres"
		dbURL = newPGURL.String()
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open connection (%s): %w", dbURL, err)
	}
	// Attempt to ping the database to ensure connection is alive
	err = db.Ping()
	if err != nil {
		_ = db.Close() // Ensure the connection is closed if not usable
		sanitizedDBURL, _ := utils.SanitizeDBURL(dbURL)
		return nil, nil, fmt.Errorf("failed to connect to database (%s): %w", sanitizedDBURL, err)
	}
	return db, func() {
		_ = db.Close()
	}, nil
}

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

func listSnapshots(db *sql.DB) (definitions.SnapshotList, error) {
	query := "SELECT datname FROM pg_database"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	list := make([]definitions.SnapshotListResult, 0)

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

		list = append(list, definitions.SnapshotListResult{
			SnapshotName: snapshotDBNameParts.SnapshotName,
			DBName:       dbName,
			CreatedAt:    snapshotDBNameParts.Timestamp,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return list, nil
}

// backupDB backs up `sourceDB` and restores it if `fn` fails
func backupDB(db *sql.DB, sourceDB string, fn func() error) error {
	if err := terminateConnections(db, sourceDB); err != nil {
		return fmt.Errorf("failed to terminate connection: %w", err)
	}
	backupDBName := "temp_emergency_backup_" + sourceDB
	if err := renameDB(db, sourceDB, backupDBName); err != nil {
		return fmt.Errorf("failed to rename snapshot: %w", err)
	}
	if err := fn(); err != nil {
		// if error, drop current source and rename backup to source
		_ = dropDB(db, sourceDB)
		_ = renameDB(db, backupDBName, sourceDB)
		return err
	}
	// is success, drop backup
	_ = dropDB(db, backupDBName)
	return nil
}

func createTemplateDB(db *sql.DB, targetDBName, sourceDBName, dbOwner string) error {
	if err := terminateConnections(db, sourceDBName); err != nil {
		return fmt.Errorf("failed to terminate connection: %w", err)
	}
	query := fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s OWNER %s;", targetDBName, sourceDBName, dbOwner)
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create template database (%s): %w", query, err)
	}
	return nil
}

func restoreDB(db *sql.DB, originalDBName, snapshotDBName, originalDBOwner string, fast bool) error {

	if fast {
		if err := dropDB(db, originalDBName); err != nil {
			return err
		}
		if err := createTemplateDB(db, originalDBName, snapshotDBName, originalDBOwner); err != nil {
			return err
		}
		return nil
	}

	return backupDB(db, originalDBName, func() error {
		if err := dropDB(db, originalDBName); err != nil {
			return err
		}
		if err := createTemplateDB(db, originalDBName, snapshotDBName, originalDBOwner); err != nil {
			return err
		}
		return nil
	})
}

func snapshotDB(db *sql.DB, originalDBName, originalDBOwner, snapshotName string) error {
	if err := terminateConnections(db, originalDBName); err != nil {
		return err
	}
	snapshotDBName, err := utils.BuildSnapshotDBName(originalDBName, snapshotName, time.Now())
	if err != nil {
		return err
	}
	return createTemplateDB(db, snapshotDBName, originalDBName, originalDBOwner)
}
