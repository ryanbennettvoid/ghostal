package mongo_db_operator

import (
	"ghostel/pkg/definitions"
	"ghostel/pkg/values"
)

type MongoDBOperatorBuilder struct{}

func (p *MongoDBOperatorBuilder) ID() string {
	return "mongo"
}

func (p *MongoDBOperatorBuilder) BuildOperator(dbURL string) (definitions.IDBOperator, error) {
	mongoURL, err := ParseMongoURL(dbURL)
	if err != nil {
		return nil, err
	}
	if mongoURL.dbURL.Scheme != "mongodb" {
		return nil, values.UnsupportedURLSchemeError
	}
	return &MongoDBOperator{
		mongoURL: mongoURL,
	}, nil
}
