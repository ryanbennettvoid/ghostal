package utils

import (
	"fmt"
	"ghostel/pkg/values"
	"strconv"
	"strings"
	"time"
)

type SnapshotDBNameParts struct {
	Name      string
	Timestamp time.Time
}

func ParseSnapshotDBName(dbName string) (SnapshotDBNameParts, error) {
	withoutPrefix := strings.TrimPrefix(dbName, values.SnapshotDBPrefix)
	partsWithoutPrefix := strings.Split(withoutPrefix, "_")
	timestamp, err := strconv.Atoi(partsWithoutPrefix[len(partsWithoutPrefix)-1])
	if err != nil {
		return SnapshotDBNameParts{}, fmt.Errorf("failed to parse database timestamp: %w", err)
	}
	name := strings.Join(partsWithoutPrefix[:len(partsWithoutPrefix)-1], "_")
	return SnapshotDBNameParts{
		Name:      name,
		Timestamp: time.UnixMilli(int64(timestamp)),
	}, nil
}
