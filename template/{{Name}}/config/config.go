package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

const assignedDefaultMsg = "Assigned default value to"

// APIConfig describes service configuration
type APIConfig struct {
	Port     string
	LogLevel string
}

// LoadConfig loads configuration
func LoadConfig() *APIConfig {
	config := &APIConfig{}
	loadEnvVars()
	validateConfig(config)
	log.Infof("%+v\n", config)

	return config
}

func loadEnvVars() *APIConfig {
	config := &APIConfig{}
	config.LogLevel = os.Getenv("LOG_LEVEL")
	config.Port = os.Getenv("PORT")
	return config
}

func validateConfig(cfg *APIConfig) {
	if cfg.Port == "" {
		cfg.Port = "3000"
		log.Infof("%s PORT: %s", assignedDefaultMsg, cfg.Port)
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
		log.Infof("%s LOG_LEVEL: %s", assignedDefaultMsg, cfg.LogLevel)
	}
}
