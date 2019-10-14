# {{Name}}

**Env vars:**  
- *LOG_LEVEL* string, required (debug, info, warning, debug, panic, fatal)
- *POSTGRES_USER* string, required
- *POSTGRES_PASSWORD* string, required
- *POSTGRES_DB* string, required
- *POSTGRES_HOST* string, required
- *POSTGRES_PORT* int, required
- *OAUTH_CLIENT_ID* string, required
- *OAUTH_CLIENT_SECRET* string, required
- *OAUTH_JWT_SECRET* string, required (random string)
- *DD_AGENT_ADDR* string, required
- *DD_SERVICE_NAME* string, required
- *DD_TRACING_ENABLED* boolean, required
- *ALLOWED_ORIGINS* string, required ("http://localhost:8081,https://ui-subscriptiontool.gismart.xyz")
- *SERVER_PORT* int, required
- *SERVER_HOST* string, required
- *POSTGRES_MAX_OPEN_CONNS* int, required
- *HTTP_TIMEOUT* int, required
- *ALLOWED_METHODS* string, required (GET,POST,PUT,DELETE,OPTIONS)
- *ALLOWED_HEADERS* string, required (Origin,X-Requested-With,Content-Type,Accept,X-Auth,Purchase-Token)
- *CREDENTIALS* boolean, required
