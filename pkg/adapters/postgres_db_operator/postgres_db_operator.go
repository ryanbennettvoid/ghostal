package postgres_db_operator

import (
	"database/sql"
	"fmt"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	"ghostel/pkg/values"
	_ "github.com/lib/pq"
)

type PostgresDBOperator struct {
	pgURL *PostgresURL
}

func CreatePostgresDBOperator(dbURL string) (*PostgresDBOperator, error) {
	pgURL, err := ParsePostgresURL(dbURL)
	if err != nil {
		return nil, err
	}
	return &PostgresDBOperator{
		pgURL: pgURL,
	}, nil
}

func (p *PostgresDBOperator) connect(useDefault bool) (*sql.DB, func(), error) {
	dbURL := p.pgURL.dbURL.String()
	if useDefault {
		newPGURL := p.pgURL.Clone()
		newPGURL.Path = "postgres"
		dbURL = newPGURL.String()
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, err
	}
	return db, func() {
		_ = db.Close()
	}, nil
}

func (p *PostgresDBOperator) checkSnapshotName(snapshotName string) error {
	if !utils.IsValidSnapshotName(snapshotName) {
		return values.NoSpecialCharsErr
	}
	list, err := p.List()
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

func (p *PostgresDBOperator) Snapshot(snapshotName string) error {
	if err := p.checkSnapshotName(snapshotName); err != nil {
		return err
	}
	db, close, err := p.connect(true)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer close()
	originalDBName := p.pgURL.DBName()
	originalDBOwner := p.pgURL.Username()
	return snapshotDB(db, originalDBName, originalDBOwner, snapshotName)
}

func (p *PostgresDBOperator) Restore(snapshotName string) error {
	db, close, err := p.connect(true)
	if err != nil {
		return err
	}
	defer close()

	list, err := listDBs(db)
	if err != nil {
		return err
	}
	for _, item := range list {
		if item.Name == snapshotName {
			originalDBName := p.pgURL.DBName()
			snapshotDBName := item.DBName
			if err := restoreDB(db, originalDBName, snapshotDBName); err != nil {
				return err
			}
			return nil
		}
	}
	return values.SnapshotNotExistsErr
}

func (p *PostgresDBOperator) Remove(snapshotName string) error {
	db, close, err := p.connect(true)
	if err != nil {
		return err
	}
	defer close()

	list, err := listDBs(db)
	if err != nil {
		return err
	}
	for _, item := range list {
		if item.Name == snapshotName {
			if err := dropDB(db, item.DBName); err != nil {
				return err
			}
			return nil
		}
	}
	return values.SnapshotNotExistsErr
}

func (p *PostgresDBOperator) List() (definitions.List, error) {
	db, close, err := p.connect(true)
	if err != nil {
		return nil, err
	}
	defer close()

	return listDBs(db)
}