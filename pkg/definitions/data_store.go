package definitions

type IDataStore interface {
	Load() ([]byte, error)
	Save([]byte) error
}
