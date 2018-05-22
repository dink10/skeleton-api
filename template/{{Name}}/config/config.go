package config

import (
	"bitbucket.org/gismart/config"
	"bitbucket.org/gismart/config/source/env"
	"bitbucket.org/gismart/config/source/file"
	"bitbucket.org/gismart/config/source/flag"
	"github.com/pkg/errors"
)

const pth = "PATH/TO/CONFIG.JSON"

func init() {
	// By default, the loader loads all the keys from the environment.
	// The loader can take other configuration source as parameters.
	err := config.NewLoader(
		// IMPORTANT: sources should be provided in accordance
		// with the priority from the minor to the major
		file.New(pth),
		env.New(),
		flag.New(),
		// Loading configuration
	).Load(&Config)

	if err != nil {
		panic(errors.Wrap(err, "error on loading config"))
	}
}
