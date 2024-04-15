package postgres_db_operator

import (
	"ghostel/pkg/definitions"
	"ghostel/pkg/values"
)

type PostgresDBOperatorBuilder struct{}

func (p *PostgresDBOperatorBuilder) ID() string {
	return "postgres"
}

func (p *PostgresDBOperatorBuilder) BuildOperator(dbURL string) (definitions.IDBOperator, error) {
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
