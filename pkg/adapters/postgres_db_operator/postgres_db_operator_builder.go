package postgres_db_operator

import (
	"ghostal/pkg/definitions"
)

type PostgresDBOperatorBuilder struct{}

func (p *PostgresDBOperatorBuilder) ID() string {
	return "Postgres"
}

func (p *PostgresDBOperatorBuilder) BuildOperator(dbURL string) (definitions.IDBOperator, error) {
	return CreatePostgresDBOperator(dbURL)
}
