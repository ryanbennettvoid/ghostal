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
	if pgURL.dbURL.Scheme != "postgresql" {
		return nil, values.UnsupportedURLSchemeError
	}
	return &PostgresDBOperator{
		pgURL: pgURL,
	}, nil
}

func (p *PostgresDBOperator) connect(useDefault bool) (*sql.DB, func(), error) {
	return createPostgresConnection(p.pgURL, useDefault)
}

func (p *PostgresDBOperator) checkSnapshotName(snapshotName string) error {
	list, err := p.ListSnapshots()
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

func (p *PostgresDBOperator) Restore(snapshotName string, fast bool) error {
	db, close, err := p.connect(true)
	if err != nil {
		return err
	}
	defer close()

	list, err := listSnapshots(db)
	if err != nil {
		return err
	}
	for _, item := range list {
		if item.SnapshotName == snapshotName {
			originalDBName := p.pgURL.DBName()
			snapshotDBName := item.DBName
			originalDBOwner := p.pgURL.Username()
			if err := restoreDB(db, originalDBName, snapshotDBName, originalDBOwner, fast); err != nil {
				return err
			}
			return nil
		}
	}
	return values.SnapshotNotExistsErr
}

func (p *PostgresDBOperator) Delete(snapshotName string) error {
	db, close, err := p.connect(true)
	if err != nil {
		return err
	}
	defer close()

	list, err := listSnapshots(db)
	if err != nil {
		return err
	}
	for _, item := range list {
		if item.SnapshotName == snapshotName {
			if err := dropDB(db, item.DBName); err != nil {
				return err
			}
			return nil
		}
	}
	return values.SnapshotNotExistsErr
}

func (p *PostgresDBOperator) ListSnapshots() (definitions.SnapshotList, error) {
	db, close, err := p.connect(true)
	if err != nil {
		return nil, err
	}
	defer close()

	return listSnapshots(db)
}
