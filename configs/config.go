package configs

import "time"

type Config struct {
	JwtKeysPath     string
	GracefulTimeout time.Duration
	Logger          LoggerConfig
}

type LoggerConfig struct {
	Level string
}
