package definitions

type ITableBuilder interface {
	BuildTable(columns []string, rows [][]string) string
}
