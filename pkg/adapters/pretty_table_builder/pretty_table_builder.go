package pretty_table_builder

import (
	"bytes"
	"github.com/olekukonko/tablewriter"
)

type PrettyTableBuilder struct{}

func NewPrettyTableBuilder() *PrettyTableBuilder {
	return &PrettyTableBuilder{}
}

func (p PrettyTableBuilder) BuildTable(columns []string, rows [][]string) string {
	buffer := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(buffer)
	table.SetRowLine(true)
	table.SetRowSeparator("-")
	table.SetHeader(columns)
	table.SetColWidth(50)
	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
	return buffer.String()
}
