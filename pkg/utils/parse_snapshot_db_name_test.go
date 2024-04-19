package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnit_ParseSnapshotDBName(t *testing.T) {
	input := "ghostalsnapshot_mydb_v1_1712976085060"
	output, err := ParseSnapshotDBName(input)
	assert.NoError(t, err)
	assert.Equal(t, "mydb", output.SourceDBName)
	assert.Equal(t, "v1", output.SnapshotName)
	assert.Equal(t, time.UnixMilli(1712976085060), output.Timestamp)
}

func TestUnit_ParseSnapshotDBNameWithUnderscores(t *testing.T) {
	input := "ghostalsnapshot_my_super_db_v2_1712976085061"
	output, err := ParseSnapshotDBName(input)
	assert.NoError(t, err)
	assert.Equal(t, "my_super_db", output.SourceDBName)
	assert.Equal(t, "v2", output.SnapshotName)
	assert.Equal(t, time.UnixMilli(1712976085061), output.Timestamp)
}

func TestUnit_ParseSnapshotDBNameCaseSensitive(t *testing.T) {
	input := "ghostalsnapshot_MyGoodDB_v3_1712976085062"
	output, err := ParseSnapshotDBName(input)
	assert.NoError(t, err)
	assert.Equal(t, "MyGoodDB", output.SourceDBName)
	assert.Equal(t, "v3", output.SnapshotName)
	assert.Equal(t, time.UnixMilli(1712976085062), output.Timestamp)
}
