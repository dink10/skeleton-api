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
	AllowedOrigins string `split_words:"true" required="true" `
	AllowedHeaders string `split_words:"true" required="true"`
	AllowedMethods string `split_words:"true" required="true"`
	Credentials    string `split_words:"true" required="true"`
	HTTPTimeout    int    `split_words:"true" required="true"`
	LogLevel       string `split_words:"true" required="true"`
	Server         struct {
		Host string `required="true"`
		Port int    `required="true"`
	}
	Postgres struct {
		Host         string `required="true"`
		Port         int    `required="true"`
		DB           string `required="true"`
		User     string `required="true"`
		Password     string `required="true"`
		MaxOpenConns int    `split_words:"true" required="true"`
	}
	Oauth    struct {
		ClientID     string `split_words:"true" required="true"`
		ClientSecret string `split_words:"true" required="true"`
		JWTSecret    string `split_words:"true" required="true"`
	}
	DD struct {
		ServiceName    string `split_words:"true" required="true"`
		AgentAddr      string `split_words:"true" required="true"`
		TracingEnabled bool   `split_words:"true" required="true"`
	}
}

func NewSchema() schema {
	return schema{}
}
