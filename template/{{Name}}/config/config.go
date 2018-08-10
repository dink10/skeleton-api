package config

import (
	"bitbucket.org/gismart/config"
	"bitbucket.org/gismart/config/source/env"
	"bitbucket.org/gismart/config/source/flag"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//const pth = "PATH/TO/CONFIG.JSON"

// Config is used to access config fields
var Config = schema{}

// InitConfig reads service configuration from a proper source
func InitConfig() {
	loadConfig(&Config)
	log.Info("Initialised config: %+v", Config)
}

func loadConfig(cfg *schema) {
	// By default, the loader loads all the keys from the environment.
	// The loader can take other configuration source as parameters.
	loader := config.NewLoader(
		// IMPORTANT: sources should be provided in accordance
		// with the priority from the minor to the major
		//file.New(pth),
		env.New(),
		flag.New(),
	)

	// Loading configuration
	err := loader.Load(cfg)
	if err != nil {
		panic(errors.Wrap(err, "error on loading config"))
	}
}
