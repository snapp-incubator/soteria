package config

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kelseyhightower/envconfig"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
)

// AppConfig is the main container of Soteria's config.
type AppConfig struct {
	AllowedAccessTypes  []string `default:"sub,pub" split_words:"true"`
	PassengerHashLength int      `split_words:"true"`
	DriverHashLength    int      `split_words:"true"`
	PassengerSalt       string   `split_words:"true"`
	DriverSalt          string   `split_words:"true"`
	Redis               *RedisConfig
	Jwt                 *JwtConfig
	Logger              *LoggerConfig
	Cache               *CacheConfig
	Mode                string `default:"debug"`
	HttpPort            int    `default:"9999" split_words:"true"`
	GrpcPort            int    `default:"50051" split_words:"true"`
	Tracer              *TracerConfig
}

// RedisConfig is all configs needed to connect to a Redis server.
type RedisConfig struct {
	Address            string        `split_words:"true"`
	Password           string        `default:"" split_words:"true"`
	ExpirationTime     int           `default:"30" split_words:"true"`
	PoolSize           int           `split_words:"true" default:"10"`
	MaxRetries         int           `split_words:"true" default:"0"`
	MinIdleConnections int           `split_words:"true" default:"5"`
	ReadTimeout        time.Duration `split_words:"true" default:"3s"`
	PoolTimeout        time.Duration `split_words:"true" default:"4s"`
	MinRetryBackoff    time.Duration `split_words:"true" default:"8ms"`
	MaxRetryBackoff    time.Duration `split_words:"true" default:"512ms"`
	IdleTimeout        time.Duration `split_words:"true" default:"300s"`
	IdleCheckFrequency time.Duration `split_words:"true" default:"60s"`
}

// CacheConfig contains configs of in memory cache.
type CacheConfig struct {
	Enabled    bool          `split_words:"true" default:"true"`
	Expiration time.Duration `split_words:"true" default:"600s"`
}

// JwtConfig contains path of the keys for JWT encryption.
type JwtConfig struct {
	KeysPath string `split_words:"true" default:"test/"`
}

// LoggerConfig is the config for logging and this kind of stuff.
type LoggerConfig struct {
	Level string `default:"warn" split_words:"true"`

	SentryEnabled bool          `default:"false" split_words:"true"`
	SentryDSN     string        `envconfig:"SENTRY_DSN"`
	SentryTimeout time.Duration `split_words:"true" default:"100ms"`
}

// TracerConfig contains all configs needed to create a tracer
type TracerConfig struct {
	Enabled      bool    `split_words:"false" default:"true"`
	ServiceName  string  `default:"soteria" split_words:"true"`
	SamplerType  string  `default:"const" split_words:"true"`
	SamplerParam float64 `default:"1" split_words:"true"`
	Host         string  `default:"localhost" split_words:"true"`
	Port         int     `default:"6831" split_words:"true"`
}

// InitConfig tries to initialize app config from env variables.
func InitConfig() AppConfig {
	appConfig := &AppConfig{}
	appConfig.Redis = &RedisConfig{}
	appConfig.Cache = &CacheConfig{}
	appConfig.Jwt = &JwtConfig{}
	appConfig.Logger = &LoggerConfig{}
	appConfig.Tracer = &TracerConfig{}

	envconfig.MustProcess("soteria", appConfig)
	envconfig.MustProcess("soteria_redis", appConfig.Redis)
	envconfig.MustProcess("soteria_cache", appConfig.Cache)
	envconfig.MustProcess("soteria_jwt", appConfig.Jwt)
	envconfig.MustProcess("soteria_logger", appConfig.Logger)
	envconfig.MustProcess("soteria_tracer", appConfig.Tracer)
	return *appConfig
}

// ReadPrivateKey will read and return private key that is used for JWT encryption
func (a *AppConfig) ReadPrivateKey(u user.Issuer) (*rsa.PrivateKey, error) {
	var fileName string
	switch u {
	case user.ThirdParty:
		fileName = fmt.Sprintf("%s%s", a.Jwt.KeysPath, "100.private.pem")
	default:
		return nil, fmt.Errorf("invalid issuer, private key not found")
	}
	pem, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// ReadPrivateKey will read and return private key that is used for JWT encryption
func (a *AppConfig) ReadPublicKey(u user.Issuer) (*rsa.PublicKey, error) {
	var fileName string
	switch u {
	case user.Driver:
		fileName = fmt.Sprintf("%s%s", a.Jwt.KeysPath, "0.pem")
	case user.Passenger:
		fileName = fmt.Sprintf("%s%s", a.Jwt.KeysPath, "1.pem")
	case user.ThirdParty:
		fileName = fmt.Sprintf("%s%s", a.Jwt.KeysPath, "100.pem")
	default:
		return nil, fmt.Errorf("invalid issuer, public key not found")
	}
	pem, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPublicKeyFromPEM(pem)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// GetAllowedAccessTypes will return all allowed access types in Soteria
func (a *AppConfig) GetAllowedAccessTypes() ([]acl.AccessType, error) {
	var allowedAccessTypes []acl.AccessType
	for _, a := range a.AllowedAccessTypes {
		at, err := toUserAccessType(a)
		if err != nil {
			return nil, fmt.Errorf("could not convert %s: %w", at, err)
		}
		allowedAccessTypes = append(allowedAccessTypes, at)
	}
	return allowedAccessTypes, nil
}

// toUserAccessType will convert string access type to it's own type
func toUserAccessType(i string) (acl.AccessType, error) {
	switch i {
	case "pub":
		return acl.Pub, nil
	case "sub":
		return acl.Sub, nil
	case "pubsub":
		return acl.PubSub, nil
	}
	return "", fmt.Errorf("%v is a invalid acces type", i)
}
