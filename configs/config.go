package configs

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type AppConfig struct {
	Redis    *RedisConfig
	Jwt      *JwtConfig
	Logger   *LoggerConfig
	HttpPort int `default:"9999" split_words:"true"`
	GrpcPort int `default:"50051" split_words:"true"`
}

type RedisConfig struct {
	Address            string        `split_words:"true"`
	Password           string        `default:"" split_words:"true"`
	ExpirationTime     int           `default:"30" split_words:"true"`
	PoolSize           int           `split_words:"true" default:"10"`
	MaxRetries         int           `split_words:"true" default:"0"`
	ReadTimeout        time.Duration `split_words:"true" default:"3s"`
	PoolTimeout        time.Duration `split_words:"true" default:"4s"`
	MinRetryBackoff    time.Duration `split_words:"true" default:"8ms"`
	MaxRetryBackoff    time.Duration `split_words:"true" default:"512ms"`
	IdleTimeout        time.Duration `split_words:"true" default:"300s"`
	IdleCheckFrequency time.Duration `split_words:"true" default:"60s"`
}

type JwtConfig struct {
	JwtKeysPath string `split_words:"true"`
}

type LoggerConfig struct {
	Level string `default:"warn" split_words:"true"`

	SentryEnabled bool          `default:"false" split_words:"true"`
	SentryDSN     string        `envconfig:"SENTRY_DSN"`
	SentryTimeout time.Duration `split_words:"true" default:"100ms"`
}

// InitConfig tries to initialize app config from env variables.
func InitConfig() AppConfig {
	appConfig := &AppConfig{}
	appConfig.Redis = &RedisConfig{}
	appConfig.Jwt = &JwtConfig{}
	appConfig.Logger = &LoggerConfig{}

	envconfig.MustProcess("soteria", appConfig)
	envconfig.MustProcess("soteria_redis", appConfig.Redis)
	envconfig.MustProcess("soteria_jwt", appConfig.Jwt)
	envconfig.MustProcess("soteria_logger", appConfig.Logger)
	return *appConfig
}
