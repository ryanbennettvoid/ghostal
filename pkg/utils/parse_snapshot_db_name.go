package utils

import (
	"fmt"
	"ghostal/pkg/values"
	"strconv"
	"strings"
	"time"
)

type SnapshotDBNameParts struct {
	SourceDBName string
	SnapshotName string
	Timestamp    time.Time
}

func ParseSnapshotDBName(fullDBName string) (SnapshotDBNameParts, error) {
	withoutPrefix := strings.TrimPrefix(fullDBName, values.SnapshotDBPrefix)

	partsAfterPrefix := strings.Split(withoutPrefix, "_")

	// go backwards from timestamp
	position := len(partsAfterPrefix) - 1

	timestamp, err := strconv.Atoi(partsAfterPrefix[position])
	if err != nil {
		return SnapshotDBNameParts{}, fmt.Errorf("failed to parse database timestamp: %w", err)
	}

	position -= 1

	snapshotName := partsAfterPrefix[position]

	position -= 1

	dbName := strings.Join(partsAfterPrefix[:position+1], "_")

	return SnapshotDBNameParts{
		SourceDBName: dbName,
		SnapshotName: snapshotName,
		Timestamp:    time.UnixMilli(int64(timestamp)),
	}, nil
}
