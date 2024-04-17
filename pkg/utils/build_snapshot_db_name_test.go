package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestUnit_BuildSnapshotDBName_ValidInput(t *testing.T) {
	inputDBName := "myDB"
	inputSnapshotName := "niceapp"
	inputTimestamp := time.Now()
	output, err := BuildSnapshotDBName(inputDBName, inputSnapshotName, inputTimestamp)
	assert.NoError(t, err)
	parts := strings.Split(output, "_")
	assert.Len(t, parts, 4)
	partsWithoutTimestamp := parts[:3]
	assert.Equal(t, []string{"ghostelsnapshot", "myDB", "niceapp"}, partsWithoutTimestamp)

	timestampStr := parts[3]
	nowUnixStr := fmt.Sprintf("%d", inputTimestamp.UnixMilli())
	assert.Equal(t, nowUnixStr, timestampStr)
}

func TestUnit_BuildSnapshotDBName_InvalidInput(t *testing.T) {
	inputDBName := "myDB"
	inputSnapshotName := "nice-app"
	inputTimestamp := time.Now()
	_, err := BuildSnapshotDBName(inputDBName, inputSnapshotName, inputTimestamp)
	assert.Error(t, err)
}
