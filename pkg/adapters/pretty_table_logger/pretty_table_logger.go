package pretty_table_logger

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

type PrettyTableLogger struct{}

func NewPrettyTableLogger() *PrettyTableLogger {
	return &PrettyTableLogger{}
}

func (p PrettyTableLogger) Log(columns []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columns)
	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
}
