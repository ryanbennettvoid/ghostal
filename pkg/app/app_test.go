package app

import (
	"context"
	"encoding/json"
	"fmt"
	"ghostel/pkg/adapters/memory_data_store"
	"ghostel/pkg/adapters/memory_logger"
	"ghostel/pkg/adapters/mongo_db_operator"
	"ghostel/pkg/adapters/postgres_db_operator"
	"ghostel/pkg/adapters/pretty_table_builder"
	"ghostel/pkg/definitions"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"strings"
	"testing"
	"time"
)

var testDBOperatorBuilders = []definitions.IDBOperatorBuilder{
	&postgres_db_operator.PostgresDBOperatorBuilder{},
	&mongo_db_operator.MongoDBOperatorBuilder{},
}

var testAppVersion = "v0.0.0"
var testLogger *memory_logger.MemoryLogger
var testTableBuilder = pretty_table_builder.NewPrettyTableBuilder()
var testDataStore *memory_data_store.MemoryDataStore

func createAndRunApp(programArgs string) error {
	testDataStore = memory_data_store.NewMemoryDataStore()
	return createAndRunAppWithDataStore(testDataStore, programArgs)
}

func createAndRunAppWithDataStore(dataStore definitions.IDataStore, programArgs string) error {
	testLogger = memory_logger.NewMemoryLogger()
	app := NewApp(testAppVersion, testDBOperatorBuilders, testLogger, testTableBuilder)
	return app.Run(dataStore, "gho", strings.Split(programArgs, " "))
}

func TestUnit_App_Version(t *testing.T) {
	assert.NoError(t, createAndRunApp("version"))
	assert.Containsf(t, testLogger.GetFullLog(), testAppVersion, "log should contain version")
}

func TestUnit_App_Help(t *testing.T) {
	assert.NoError(t, createAndRunApp("help"))
	fullLog := testLogger.GetFullLog()
	assert.Contains(t, fullLog, "gho version")
	assert.Contains(t, fullLog, "gho help")
	assert.Contains(t, fullLog, "gho init <project_name> <database_name>")
	assert.Contains(t, fullLog, "gho select <project_name>")
	assert.Contains(t, fullLog, "gho set <key> <value>")
	assert.Contains(t, fullLog, "gho status")
	assert.Contains(t, fullLog, "gho snapshot <snapshot_name>")
	assert.Contains(t, fullLog, "gho restore <snapshot_name>")
	assert.Contains(t, fullLog, "gho rm <snapshot_name>")
	assert.Contains(t, fullLog, "gho ls")
}

func TestUnit_App_Init(t *testing.T) {
	assert.NoError(t, createAndRunApp("init xxx postgresql://localhost"))
	var c definitions.ConfigData
	assert.NoError(t, json.Unmarshal(testDataStore.Data, &c), "should have saved correct config data")
	assert.Equal(t, "xxx", c.SelectedProject, "config data selected project should be set")
}

func TestUnit_App_Init_Errors(t *testing.T) {
	assert.Errorf(t, createAndRunApp("init"), "should fail if no options provided")
	assert.Errorf(t, createAndRunApp("init xxx"), "should fail if no database URL provided")
	assert.Errorf(t, createAndRunApp("init xxx invalidUrl"), "should fail if database url is invalid")
}

func TestUnit_App_Select(t *testing.T) {
	dataStore := memory_data_store.NewMemoryDataStore()
	configDataToSave, err := json.Marshal(definitions.ConfigData{
		SelectedProject: "aaa",
		Projects: []definitions.Project{
			{
				Name:      "aaa",
				DBURL:     "postgresql://localhost",
				CreatedAt: time.Time{},
			},
			{
				Name:      "bbb",
				DBURL:     "mongodb://localhost",
				CreatedAt: time.Time{},
			},
		},
	})
	assert.NoError(t, err)
	dataStore.Data = configDataToSave
	assert.NoError(t, createAndRunAppWithDataStore(dataStore, "select bbb"))
	var c definitions.ConfigData
	assert.NoError(t, json.Unmarshal(dataStore.Data, &c), "should have saved correct config data")
	assert.Equal(t, "bbb", c.SelectedProject)
}

func TestUnit_App_Set(t *testing.T) {
	dataStore := memory_data_store.NewMemoryDataStore()
	configDataToSave, err := json.Marshal(definitions.ConfigData{
		SelectedProject: "aaa",
		Projects: []definitions.Project{
			{
				Name:      "aaa",
				DBURL:     "postgresql://localhost",
				CreatedAt: time.Time{},
			},
		},
	})
	assert.NoError(t, err)
	dataStore.Data = configDataToSave
	{
		assert.NoError(t, createAndRunAppWithDataStore(dataStore, "set fastRestore true"))
		var c definitions.ConfigData
		assert.NoError(t, json.Unmarshal(dataStore.Data, &c), "should have saved correct config data")
		assert.True(t, *c.Projects[0].FastRestore)
	}
	{
		assert.NoError(t, createAndRunAppWithDataStore(dataStore, "set fastRestore false"))
		var c definitions.ConfigData
		assert.NoError(t, json.Unmarshal(dataStore.Data, &c), "should have saved correct config data")
		assert.False(t, *c.Projects[0].FastRestore)
	}
	{
		assert.Error(t, createAndRunAppWithDataStore(dataStore, "set fastRestore xxx"))
	}
}

func TestUnit_App_Status(t *testing.T) {
	dataStore := memory_data_store.NewMemoryDataStore()
	configDataToSave, err := json.Marshal(definitions.ConfigData{
		SelectedProject: "aaa",
		Projects: []definitions.Project{
			{
				Name:      "aaa",
				DBURL:     "postgresql://localhost",
				CreatedAt: time.Time{},
			},
			{
				Name:      "bbb",
				DBURL:     "mongodb://localhost",
				CreatedAt: time.Time{},
			},
		},
	})
	assert.NoError(t, err)
	dataStore.Data = configDataToSave
	assert.NoError(t, createAndRunAppWithDataStore(dataStore, "status"))
	fullLog := testLogger.GetFullLog()
	assert.Contains(t, fullLog, "aaa")
	assert.Contains(t, fullLog, "postgresql://localhost")
	assert.Contains(t, fullLog, "bbb")
	assert.Contains(t, fullLog, "mongodb://localhost")
}

// ---------

func createPostgresContainer() (string, func()) {
	dbName := "mydb"
	username := "pguser"
	password := "pgpass"
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15.1-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     username,
			"POSTGRES_PASSWORD": password,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to start container: %s", err))
	}

	host, err := container.Host(context.Background())
	if err != nil {
		panic(err)
	}

	mappedPort, err := container.MappedPort(context.Background(), "5432")
	if err != nil {
		panic(err)
	}

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, mappedPort.Port(), dbName)
	return dbURL, func() {
		_ = container.Terminate(ctx)
	}
}

func createMongoContainer() (string, func()) {
	dbName := "mydb"
	username := "mongouser"
	password := "mongopass"
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:7.0.5",
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_DATABASE":      dbName,
			"MONGO_INITDB_ROOT_USERNAME": username,
			"MONGO_INITDB_ROOT_PASSWORD": password,
		},
		WaitingFor: wait.ForListeningPort("27017/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to start container: %s", err))
	}

	host, err := container.Host(context.Background())
	if err != nil {
		panic(err)
	}

	mappedPort, err := container.MappedPort(context.Background(), "27017")
	if err != nil {
		panic(err)
	}

	dbURL := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?tls=false", username, password, host, mappedPort.Port(), dbName)
	return dbURL, func() {
		_ = container.Terminate(ctx)
	}
}

func TestIntegration_App_SnapshotSmokePostgres(t *testing.T) {
	dbURL, cleanup := createPostgresContainer()
	defer cleanup()
	postgres_db_operator.WritePostgresSeedData(dbURL)
	runSmokeTest(t, "pg_local", dbURL)
}

func TestIntegration_App_SnapshotSmokeMongo(t *testing.T) {
	dbURL, cleanup := createMongoContainer()
	defer cleanup()
	mongo_db_operator.WriteMongoDBSeedData(dbURL)
	runSmokeTest(t, "mongo_local", dbURL)
}

func assertLogContains(t *testing.T, substring string, positive bool, fn func()) {
	// reset memory logger before listing snapshots
	testLogger = memory_logger.NewMemoryLogger()
	fn()
	fullLog := testLogger.GetFullLog()
	if positive {
		assert.Containsf(t, fullLog, substring, fmt.Sprintf("log should contain \"%s\"", substring))
	} else {
		assert.NotContainsf(t, fullLog, substring, fmt.Sprintf("log should NOT contain \"%s\"", substring))
	}
}

func runSmokeTest(t *testing.T, projectName, dbURL string) {
	dataStore := memory_data_store.NewMemoryDataStore()
	initCmd := fmt.Sprintf("init %s %s", projectName, dbURL)

	assert.NoErrorf(t, createAndRunAppWithDataStore(dataStore, initCmd), "should initialize project")
	assert.NoErrorf(t, createAndRunAppWithDataStore(dataStore, "snapshot v1"), "should create snapshot")
	assert.Errorf(t, createAndRunAppWithDataStore(dataStore, "snapshot v1"), "should fail to create snapshot with duplicate name")

	assertLogContains(t, "v1", true, func() {
		assert.NoErrorf(t, createAndRunAppWithDataStore(dataStore, "ls"), "should list snapshots")
	})

	assert.Errorf(t, createAndRunAppWithDataStore(dataStore, "restore v2"), "should fail to restore non-existant snapshot")

	assert.NoErrorf(t, createAndRunAppWithDataStore(dataStore, "restore v1"), "should restore snapshot")
	assert.NoErrorf(t, createAndRunAppWithDataStore(dataStore, "restore v1"), "should restore snapshot again")

	assertLogContains(t, "v1", true, func() {
		assert.NoErrorf(t, createAndRunAppWithDataStore(dataStore, "ls"), "should list snapshots")
	})

	// remove snapshot
	assert.NoErrorf(t, createAndRunAppWithDataStore(dataStore, "rm v1"), "should remove snapshot")

	assertLogContains(t, "v1", false, func() {
		assert.NoErrorf(t, createAndRunAppWithDataStore(dataStore, "ls"), "should list snapshots")
	})

	assert.NoErrorf(t, createAndRunAppWithDataStore(dataStore, "snapshot xxx"), "should create snapshot")
	assert.NoErrorf(t, createAndRunAppWithDataStore(dataStore, "set fastRestore true"), "should set fastRestore to true for project")
	assert.NoErrorf(t, createAndRunAppWithDataStore(dataStore, "restore xxx"), "should restore snapshot (fast)")
}
