package definitions

type ITableLogger interface {
	Log(columns []string, rows [][]string)
}
