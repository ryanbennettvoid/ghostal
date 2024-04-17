package definitions

import (
	"ghostel/pkg/utils"
	"time"
)

type SnapshotListResult struct {
	SnapshotName string
	DBName       string
	CreatedAt    time.Time
}

type SnapshotList []SnapshotListResult

func (list SnapshotList) TableInfo() ([]string, [][]string) {
	columns := []string{"Name", "Created", "Timestamp"}
	rows := make([][]string, len(list))
	for idx := range list {
		item := list[idx]
		relativeTime := utils.ToRelativeTime(item.CreatedAt, time.Now())
		formattedTime := item.CreatedAt.Format("2006-01-02 15:04:05")
		rows[idx] = []string{item.SnapshotName, relativeTime, formattedTime}
	}
	return columns, rows
}

type IDBOperator interface {
	Snapshot(snapshotName string) error
	Restore(snapshotName string, fast bool) error
	Delete(snapshotName string) error
	ListSnapshots() (SnapshotList, error)
}
