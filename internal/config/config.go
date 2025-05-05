package config

import (
	"os"

	flag "github.com/spf13/pflag"
)

type Configuration struct {
	Environment *string
	HTTPPort    *string
}

var (
	Config *Configuration

	env = flag.String(
		"env",
		"production",
		"Environment")

	httpPort = flag.String(
		"http-port",
		"5001",
		"The port to serve on")
)

func updateStringEnvVariable(defValue *string, key string) *string {
	val := os.Getenv(key)

	if val == "" {
		return defValue
	}

	return &val
}

func init() {
	flag.Parse()

	env = updateStringEnvVariable(env, "ENVIRONMENT")
	httpPort = updateStringEnvVariable(httpPort, "HTTP_PORT")

	Config = &Configuration{
		Environment: env,
		HTTPPort:    httpPort,
	}
}
