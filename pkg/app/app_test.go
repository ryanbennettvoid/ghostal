package app

import (
	"ghostel/pkg/adapters/memory_data_store"
	"ghostel/pkg/adapters/memory_logger"
	"ghostel/pkg/adapters/pretty_table_builder"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var testAppVersion = "v0.0.0"
var testLogger *memory_logger.MemoryLogger
var testTableBuilder = pretty_table_builder.NewPrettyTableBuilder()

func createAndRunApp(programArgs string) error {
	app := NewApp(testAppVersion, testLogger, testTableBuilder)
	dataStore := memory_data_store.NewMemoryDataStore()
	return app.Run(dataStore, "gho", strings.Split(programArgs, " "))
}

func TestUnit_App_Version(t *testing.T) {
	testLogger = memory_logger.NewMemoryLogger()
	assert.NoError(t, createAndRunApp("version"))
	assert.Containsf(t, testLogger.GetFullLog(), testAppVersion, "log should contain version")
}

func TestUnit_App_Help(t *testing.T) {
	testLogger = memory_logger.NewMemoryLogger()
	assert.NoError(t, createAndRunApp("help"))
	fullLog := testLogger.GetFullLog()
	assert.Contains(t, fullLog, "gho version")
	assert.Contains(t, fullLog, "gho help")
	assert.Contains(t, fullLog, "gho init <project_name> <database_name>")
	assert.Contains(t, fullLog, "gho select <project_name>")
	assert.Contains(t, fullLog, "gho status")
	assert.Contains(t, fullLog, "gho snapshot <snapshot_name>")
	assert.Contains(t, fullLog, "gho restore <snapshot_name>")
	assert.Contains(t, fullLog, "gho rm <snapshot_name>")
	assert.Contains(t, fullLog, "gho ls")
}

//func TestUnit_App_Init(t *testing.T) {
//	testLogger = memory_logger.NewMemoryLogger()
//	assert.NoError(t, createAndRunApp("init xxx yyy"))
//}
