package main

import (
	"ghostel/pkg/adapters/logrus_logger"
	"ghostel/pkg/adapters/pretty_table_logger"
	"ghostel/pkg/app"
	"os"
	"time"
)

var (
	Version string
)

var logger = logrus_logger.NewLogrusLogger()
var tableLogger = pretty_table_logger.NewPrettyTableLogger()

var start = time.Now()

func exit(err error) {
	if err == nil {
		logger.Info("Done in %.3fs.\n", time.Now().Sub(start).Seconds())
		os.Exit(0)
	} else {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func main() {
	executable := os.Args[0]
	args := os.Args[1:]
	app := app.NewApp(Version, logger, tableLogger)
	exit(app.Run(executable, args))
}
