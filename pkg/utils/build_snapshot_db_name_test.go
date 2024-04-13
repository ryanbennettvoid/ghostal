package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestUnit_BuildSnapshotDBName_ValidInput(t *testing.T) {
	inputName := "niceapp"
	inputTimestamp := time.Now()
	output := BuildSnapshotDBName(inputName, inputTimestamp)
	parts := strings.Split(output, "_")
	assert.Len(t, parts, 4)
	partsWithoutTimestamp := parts[:3]
	assert.Equal(t, []string{"ghostel", "snapshot", "niceapp"}, partsWithoutTimestamp)

	timestampStr := parts[3]
	nowUnixStr := fmt.Sprintf("%d", inputTimestamp.UnixMilli())
	assert.Equal(t, nowUnixStr, timestampStr)
}
