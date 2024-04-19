package mongo_db_operator

import (
	"ghostal/pkg/definitions"
)

type MongoDBOperatorBuilder struct{}

func (p *MongoDBOperatorBuilder) ID() string {
	return "MongoDB"
}

func (p *MongoDBOperatorBuilder) BuildOperator(dbURL string) (definitions.IDBOperator, error) {
	return CreateMongoDBOperator(dbURL)
}
