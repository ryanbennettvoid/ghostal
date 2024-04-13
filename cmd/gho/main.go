package main

import (
	"ghostel/pkg/app"
	"os"
	"time"
)

var logger = app.GetGlobalLogger()

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
	app := app.NewApp(logger)
	exit(app.Run(os.Args[1:]))
}
