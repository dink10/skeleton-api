package config

import (
	"github.com/koding/multiconfig"
	log "github.com/sirupsen/logrus"
)

// TemplateConfig TODO: Change it's name and fields according to your requirements!
type TemplateConfig struct {
	Port     string `default:"8282"`
	LogLevel string `default:"debug"`
}

// LoadConfig loads configuration
func LoadConfig() *TemplateConfig {
	config := &TemplateConfig{}
	m := multiconfig.New()
	log.Infof("Loading configuration...")
	err := m.Load(config)
	if err != nil {
		log.Fatalf("Failed to load configuration. %v", err)
	}

	err = m.Validate(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("%+v\n", config)

	return config
}
