package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestBuildSnapshotDBName_ValidInput(t *testing.T) {
	input := "niceapp"
	output := BuildSnapshotDBName(input)
	parts := strings.Split(output, "_")
	assert.Len(t, parts, 4)
	partsWithoutTimestamp := parts[:3]
	assert.Equal(t, []string{"ghostel", "snapshot", "niceapp"}, partsWithoutTimestamp)

	// NOTE: comparing timestamps has a non-zero chance of random failure
	timestampStr := parts[3]
	nowUnixStr := fmt.Sprintf("%d", time.Now().Unix())
	assert.Equalf(t, nowUnixStr[:9], timestampStr[:9], "first 9 digits of timestamp match current time")
}
