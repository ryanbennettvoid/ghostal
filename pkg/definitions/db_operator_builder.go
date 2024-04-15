package definitions

type IDBOperatorBuilder interface {
	ID() string
	BuildOperator(dbURL string) (IDBOperator, error)
}
