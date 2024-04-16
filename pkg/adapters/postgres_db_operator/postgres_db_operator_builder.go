package postgres_db_operator

import (
	"ghostel/pkg/definitions"
)

type PostgresDBOperatorBuilder struct{}

func (p *PostgresDBOperatorBuilder) ID() string {
	return "postgres"
}

func (p *PostgresDBOperatorBuilder) BuildOperator(dbURL string) (definitions.IDBOperator, error) {
	return CreatePostgresDBOperator(dbURL)
}
