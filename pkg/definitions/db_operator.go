package definitions

import "time"

type ListResult struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type IDBOperator interface {
	Snapshot(snapshotName string) error
	Restore(snapshotName string) error
	Remove(snapshotName string) error
	List() ([]ListResult, error)
}
