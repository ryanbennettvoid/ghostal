package definitions

import (
	"ghostel/pkg/utils"
	"time"
)

type ListResult struct {
	Name      string    `json:"name"`
	DBName    string    `json:"db_name"`
	CreatedAt time.Time `json:"created_at"`
}

type List []ListResult

func (list List) Print(logger ITableLogger) {
	columns := []string{"Name", "Created", "Timestamp"}
	rows := make([][]string, len(list))
	for idx := range list {
		item := list[idx]
		relativeTime := utils.ToRelativeTime(item.CreatedAt, time.Now())
		formattedTime := item.CreatedAt.Format("2006-01-02 15:04:05")
		rows[idx] = []string{item.Name, relativeTime, formattedTime}
	}
	logger.Log(columns, rows)
}

type IDBOperator interface {
	Snapshot(snapshotName string) error
	Restore(snapshotName string) error
	Delete(snapshotName string) error
	List() (List, error)
}
