package config

// Scheme describes service configuration for JSON
// Available tag params:
// default - optional argument that sets a default value for field when it's not found any value
// alias - optional argument that sets additional another alias for field, commonly for environment variables names
// source - optional argument that sets a configuration source name directly
// required - optional argument that indicates required field
// See an example:
// type scheme struct {
//    Server struct {
//        Host string `config:"srv-host, alias=SRV_HOST, source=env, required"`
//        Port int    `config:"srv-port, default=8080"`
//  }

type schema struct {
	Server struct {
		Host string `config:"srv-host"`
		Port int    `config:"srv-port"`
	}
	Database struct {
		Host     string `config:"db-host"`
		Port     int    `config:"db-port"`
		Name     string `config:"db-name"`
		User     string `config:"db-user"`
		Password string `config:"db-password"`
	}
	Logger struct {
		LogLevel         string `config:"loglevel"`
		SentryDSN        string `config:"sentry-dsn"`
		DataDogEnv       string `config:"datadog-env"`
		DataDogAgentAddr string `config:"datadog-agent-addr"`
	}
}
