package utils

import (
	"fmt"
	"ghostel/pkg/values"
	"time"
)

func BuildFullSnapshotName(snapshotName string) string {
	return fmt.Sprintf("%s%s_%d", values.SnapshotDBPrefix, snapshotName, time.Now().UnixMilli())
}
