package utils

import (
	"fmt"
	"ghostel/pkg/values"
	"time"
)

func BuildSnapshotDBName(snapshotName string, timestamp time.Time) string {
	return fmt.Sprintf("%s%s_%d", values.SnapshotDBPrefix, snapshotName, timestamp.UnixMilli())
}
