package main

import (
	"ghostal/pkg/adapters/file_data_store"
	"ghostal/pkg/adapters/logrus_logger"
	"ghostal/pkg/adapters/mongo_db_operator"
	"ghostal/pkg/adapters/postgres_db_operator"
	"ghostal/pkg/adapters/pretty_table_builder"
	"ghostal/pkg/app"
	"ghostal/pkg/definitions"
	"ghostal/pkg/values"
	"os"
	"time"
)

var (
	Version string
)

var dbOperatorBuilders = []definitions.IDBOperatorBuilder{
	&postgres_db_operator.PostgresDBOperatorBuilder{},
	&mongo_db_operator.MongoDBOperatorBuilder{},
}

var logger = logrus_logger.NewLogrusLogger()
var tableBuilder = pretty_table_builder.NewPrettyTableBuilder()

var start = time.Now()

func exit(err error) {
	if err == nil {
		logger.Passthrough("Done in %.3fs.\n", time.Since(start).Seconds())
		os.Exit(0)
	} else {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func main() {
	executable := os.Args[0]
	args := os.Args[1:]
	app := app.NewApp(Version, dbOperatorBuilders, logger, tableBuilder)
	dataStore := file_data_store.NewFileDataStore(values.DefaultConfigFilepath)
	exit(app.Run(dataStore, executable, args))
}
