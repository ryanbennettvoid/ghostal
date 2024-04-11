package postgres_db_adapter

import (
	"database/sql"
	"errors"
	"fmt"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	"ghostel/pkg/values"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

var NoSpecialCharsErr = errors.New("no special characters allowed in snapshot name")

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

func (p *PostgresDBAdapter) Snapshot(snapshotName string) error {
	if !utils.IsAlphanumeric(snapshotName) {
		return NoSpecialCharsErr
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
	if !utils.IsAlphanumeric(snapshotName) {
		return NoSpecialCharsErr
	}
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDBAdapter) Remove(snapshotName string) error {
	if !utils.IsAlphanumeric(snapshotName) {
		return NoSpecialCharsErr
	}
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDBAdapter) List() ([]definitions.ListResult, error) {
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

		nameAndTimestamp := strings.Split(dbName, values.SnapshotDBPrefix)[1]
		parts := strings.Split(nameAndTimestamp, "_")
		name := parts[0]
		timestamp, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse database timestamp: %w", err)
		}

		list = append(list, definitions.ListResult{
			Name:      name,
			CreatedAt: time.Unix(int64(timestamp/1000), 0),
		})
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Fatal("Error during rows iteration: ", err)
	}

	return list, err
}
