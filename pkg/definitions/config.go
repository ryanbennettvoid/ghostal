package definitions

type IConfig interface {
	Get(key string) (string, error)
	Set(key, value string) error
}
