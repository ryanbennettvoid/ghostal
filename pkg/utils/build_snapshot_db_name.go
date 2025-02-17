package utils

import (
	"fmt"
	"ghostal/pkg/values"
	"time"
	"unicode"
)

func BuildSnapshotDBName(sourceDBName, snapshotName string, timestamp time.Time) (string, error) {
	// allow alphanumeric chars and underscore
	for _, c := range snapshotName {
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_') {
			return "", values.NoSpecialCharsErr
		}
	}
	result := fmt.Sprintf("%s%s_%s_%d", values.SnapshotDBPrefix, sourceDBName, snapshotName, timestamp.UnixMilli())
	return result, nil
}
