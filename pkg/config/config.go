package config

import "os"

const (
	standardAllowOrigin        = "localhost:3000"
	standardSchema             = "http://"
	standardPort               = "8080"
	standardURLDataBase        = "postgres://postgres:postgres@localhost:5432/marketplace?sslmode=disable"
	standardPathToRoot         = "."
	standardOutputLogPath      = "stdout /var/log/backend/logs.json"
	standardErrorOutputLogPath = "stderr /var/log/backend/err_logs.json"

	envAllowOrigin        = "ALLOW_ORIGIN"
	envSchema             = "SCHEMA"
	envPortBackend        = "PORT_BACKEND"
	envURLDataBase        = "URL_DATA_BASE"
	envPathToRoot         = "PATH_TO_ROOT"
	envOutputLogPath      = "OUTPUT_LOG_PATH"
	envErrorOutputLogPath = "ERROR_OUTPUT_LOG_PATH"
)

type Config struct {
	AllowOrigin        string
	Schema             string
	PortServer         string
	URLDataBase        string
	PathToRoot         string
	OutputLogPath      string
	ErrorOutputLogPath string
}

func New() *Config {
	return &Config{
		AllowOrigin:        getEnvStr(envAllowOrigin, standardAllowOrigin),
		Schema:             getEnvStr(envSchema, standardSchema),
		PortServer:         getEnvStr(envPortBackend, standardPort),
		URLDataBase:        getEnvStr(envURLDataBase, standardURLDataBase),
		PathToRoot:         getEnvStr(envPathToRoot, standardPathToRoot),
		OutputLogPath:      getEnvStr(envOutputLogPath, standardOutputLogPath),
		ErrorOutputLogPath: getEnvStr(envErrorOutputLogPath, standardErrorOutputLogPath),
	}
}

func getEnvStr(name string, defaultValue string) string {
	result, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return result
}
