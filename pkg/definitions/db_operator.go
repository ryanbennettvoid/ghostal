package definitions

import (
	"ghostel/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"os"
	"time"
)

type ListResult struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type List []ListResult

func (list List) Print() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Created", "Timestamp"})
	for _, item := range list {
		relativeTime := utils.ToRelativeTime(item.CreatedAt)
		formattedTime := item.CreatedAt.Format("2006-01-02 15:04:05")
		table.Append([]string{item.Name, relativeTime, formattedTime})
	}
	table.Render()
}

type IDBOperator interface {
	Snapshot(snapshotName string) error
	Restore(snapshotName string) error
	Remove(snapshotName string) error
	List() (List, error)
}
