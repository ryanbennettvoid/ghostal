package postgres_db_adapter

import (
	"database/sql"
	"fmt"
	"ghostel/pkg/definitions"
	_ "github.com/lib/pq"
	"time"
)

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
	db, close, err := p.connect()
	if err != nil {
		return err
	}
	defer close()
	fullSnapshotName := fmt.Sprintf("%s_%d", snapshotName, time.Now().Unix())
	originalDBName := p.pgURL.DBName()
	owner := p.pgURL.Username()
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s OWNER %s;", fullSnapshotName, originalDBName, owner))
	return err
}

func (p *PostgresDBAdapter) Restore(snapshotName string) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDBAdapter) Remove(snapshotName string) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDBAdapter) List() ([]definitions.ListResult, error) {
	//TODO implement me
	panic("implement me")
}
