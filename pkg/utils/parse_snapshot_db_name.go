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

	lastUnderscore := strings.LastIndex(withoutPrefix, "_")
	if lastUnderscore == -1 {
		return SnapshotDBNameParts{}, fmt.Errorf("invalid snapshot DB name: missing timestamp separator")
	}
	timestampStr := withoutPrefix[lastUnderscore+1:]
	tsInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return SnapshotDBNameParts{}, fmt.Errorf("failed to parse database timestamp: %w", err)
	}

	withoutTimestamp := withoutPrefix[:lastUnderscore]

	secondLastUnderscore := strings.LastIndex(withoutTimestamp, "_")
	if secondLastUnderscore == -1 {
		return SnapshotDBNameParts{}, fmt.Errorf("invalid snapshot DB name: missing snapshot name separator")
	}
	snapshotName := withoutTimestamp[secondLastUnderscore+1:]
	sourceDBName := withoutTimestamp[:secondLastUnderscore]

	return SnapshotDBNameParts{
		SourceDBName: sourceDBName,
		SnapshotName: snapshotName,
		Timestamp:    time.UnixMilli(tsInt),
	}, nil
}
