package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnit_ParseSnapshotDBName(t *testing.T) {
	input := "ghostel_snapshot_mydb_1712976085060"
	output, err := ParseSnapshotDBName(input)
	assert.NoError(t, err)
	assert.Equal(t, "mydb", output.Name)
	assert.Equal(t, time.UnixMilli(1712976085060), output.Timestamp)
}
