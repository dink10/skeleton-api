package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Config is used to access config fields
var Config = NewSchema()

func Init(cfg interface{}) {
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Error(errors.Wrap(err, "error on loading config"))
	}
	log.Info("Initialised config")
}
