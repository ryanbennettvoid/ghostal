package utils

import (
	"fmt"
	"ghostal/pkg/values"
	"time"
	"unicode"
)

func BuildSnapshotDBName(sourceDBName, snapshotName string, timestamp time.Time) (string, error) {
	// only allow alphanumeric chars
	for _, c := range snapshotName {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return "", values.NoSpecialCharsErr
		}
	}
	result := fmt.Sprintf("%s%s_%s_%d", values.SnapshotDBPrefix, sourceDBName, snapshotName, timestamp.UnixMilli())
	return result, nil
}
