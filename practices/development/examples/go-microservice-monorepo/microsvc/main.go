package main

import (
	"example.com/sample/commonlib/logger"
	"example.com/sample/microsvc/options"
	_ "go.uber.org/mock/mockgen/model"
)

// TODO update the swagger documentation header here to be customized for your microservice

//go:generate swag init --parseVendor --parseDependency
// @title Microservice Template (Sample)
// @version 1.0
// @description This microservice template demonstrates how to write the code for a Go microservice.

// @contact.name   API Support
// @contact.url    https://example.com/contact-us

// @host      127.0.0.1
// @BasePath  /

// main is the entrypoint of the microservice
func main() {
	db := PrepareSubsystems()

	logger.Log.Info("Starting example microservice...")
	router := Bootstrap(db)
	router.Listen(options.Registry)
}
