package postgres_db_adapter

import (
	"database/sql"
	"errors"
	"fmt"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	"ghostel/pkg/values"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

var NoSpecialCharsErr = errors.New("snapshot name can only contain alphanumeric characters or underscores")
var SnapshotNameTakenErr = errors.New("snapshot name already used")
var SnapshotNotExistsErr = errors.New("snapshot does not exist")

type PostgresDBAdapter struct {
	pgURL *PostgresURL
}

func CreatePostgresDBAdapter(dbURL string) (*PostgresDBAdapter, error) {
	pgURL, err := ParsePostgresURL(dbURL)
	if err != nil {
		return nil, err
	}
	return &PostgresDBAdapter{
		pgURL: pgURL,
	}, nil
}

func (p *PostgresDBAdapter) connect() (*sql.DB, func(), error) {
	db, err := sql.Open("postgres", p.pgURL.dbURL.String())
	if err != nil {
		return nil, nil, err
	}
	return db, func() {
		_ = db.Close()
	}, nil
}

func (p *PostgresDBAdapter) GetScheme() string {
	return p.pgURL.Scheme()
}

func (p *PostgresDBAdapter) checkSnapshotName(snapshotName string) error {
	if !utils.IsValidSnapshotName(snapshotName) {
		return NoSpecialCharsErr
	}
	list, err := p.List()
	if err != nil {
		return err
	}
	for _, item := range list {
		if item.Name == snapshotName {
			return SnapshotNameTakenErr
		}
	}
	return nil
}

func (p *PostgresDBAdapter) Snapshot(snapshotName string) error {
	if err := p.checkSnapshotName(snapshotName); err != nil {
		return err
	}
	db, close, err := p.connect()
	if err != nil {
		return err
	}
	defer close()
	fullSnapshotName := fmt.Sprintf("%s%s_%d", values.SnapshotDBPrefix, snapshotName, time.Now().UnixMilli())
	originalDBName := p.pgURL.DBName()
	owner := p.pgURL.Username()
	query := fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s OWNER %s;", fullSnapshotName, originalDBName, owner)
	_, err = db.Exec(query)
	return err
}

func (p *PostgresDBAdapter) Restore(snapshotName string) error {
	if err := p.checkSnapshotName(snapshotName); err != nil {
		return err
	}
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDBAdapter) removeDb(dbName string) error {
	db, close, err := p.connect()
	if err != nil {
		return err
	}
	defer close()

	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", pq.QuoteIdentifier(dbName))
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDBAdapter) Remove(snapshotName string) error {
	list, err := p.List()
	if err != nil {
		return err
	}
	for _, item := range list {
		if item.Name == snapshotName {
			if err := p.removeDb(item.DBName); err != nil {
				return err
			}
			return nil
		}
	}
	return SnapshotNotExistsErr
}

func (p *PostgresDBAdapter) List() (definitions.List, error) {
	db, close, err := p.connect()
	if err != nil {
		return nil, err
	}
	defer close()

	query := "SELECT datname FROM pg_database"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	list := make([]definitions.ListResult, 0)
	// Iterate through the result set
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if !strings.HasPrefix(dbName, values.SnapshotDBPrefix) {
			continue
		}

		withoutPrefix := strings.TrimPrefix(dbName, values.SnapshotDBPrefix)
		partsWithoutPrefix := strings.Split(withoutPrefix, "_")
		timestamp, err := strconv.Atoi(partsWithoutPrefix[len(partsWithoutPrefix)-1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse database timestamp: %w", err)
		}
		name := strings.Join(partsWithoutPrefix[:len(partsWithoutPrefix)-1], "_") // without suffix

		list = append(list, definitions.ListResult{
			Name:      name,
			DBName:    dbName,
			CreatedAt: time.Unix(int64(timestamp/1000), 0),
		})
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Fatal("Error during rows iteration: ", err)
	}

	return list, err
}
