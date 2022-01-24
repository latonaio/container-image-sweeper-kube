package main

import (
	"github.com/latonaio/container-image-sweeper-kube/cmd/app"
	"github.com/latonaio/golang-logging-library/logger"
)

var log = logger.NewLogger()

func main() {
	cmd := app.Command()
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
