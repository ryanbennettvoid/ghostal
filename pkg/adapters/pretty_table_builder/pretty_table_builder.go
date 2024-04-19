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
	table.SetAutoWrapText(true)
	table.SetHeader(columns)
	table.SetHeaderLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetColWidth(50)
	table.SetBorder(false)
	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
	return buffer.String()
}
