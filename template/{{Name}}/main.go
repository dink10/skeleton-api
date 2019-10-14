package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/gismart/ddtracer"

	"bitbucket.org/gismart/{{Name}}/config"
	"bitbucket.org/gismart/{{Name}}/database"
	"bitbucket.org/gismart/{{Name}}/server"
	"bitbucket.org/gismart/{{Name}}/services/logger"
)

// @title Subscription Tool
func main() {
	config.Init(&config.Config)
	logger.Init(config.Config.LogLevel)

	err := database.PostgresConnect()
	if err != nil {
		log.Fatalf("Unable to initialize postgres: %s", err.Error())
	}

	if config.Config.DD.TracingEnabled {
		ddtracer.InitHttpTracer(config.Config.DD.AgentAddr, config.Config.DD.ServiceName)
		defer ddtracer.Stop()
	}

	defer func() {
		err := database.PostgresClose()
		if err != nil {
			log.Error(err)
		}
	}()

	addr := fmt.Sprintf(":%d", config.Config.Server.Port)
	server.Run(addr, config.Config.DD.TracingEnabled)
}
