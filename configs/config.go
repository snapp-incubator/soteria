package configs

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kelseyhightower/envconfig"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"io/ioutil"
	"log"
	"time"
)

type AppConfig struct {
	AllowedAccessTypes  []string `default:"sub,pub" split_words:"true"`
	PassengerHashLength int      `split_words:"true"`
	DriverHashLength    int      `split_words:"true"`
	PassengerSalt       string   `split_words:"true"`
	DriverSalt          string   `split_words:"true"`
	Redis               *RedisConfig
	Jwt                 *JwtConfig
	Logger              *LoggerConfig
	HttpPort            int `default:"9999" split_words:"true"`
	GrpcPort            int `default:"50051" split_words:"true"`
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

func (a *AppConfig) ReadPrivateKey(u string) (*rsa.PrivateKey, error) {
	var fileName string
	switch u {
	case user.ThirdParty:
		fileName = fmt.Sprintf("%s%s", a.Jwt.JwtKeysPath, "100.private.pem")
	default:
		return nil, fmt.Errorf("invalid user, private key not found")
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

func (a *AppConfig) GetAllowedAccessTypes() []user.AccessType {
	var allowedAccessTypes []user.AccessType
	for _, a := range a.AllowedAccessTypes {
		at, err := toUserAccessType(a)
		if err != nil {
			log.Fatal(err)
		}
		allowedAccessTypes = append(allowedAccessTypes, at)
	}
	return allowedAccessTypes
}

func toUserAccessType(i string) (user.AccessType, error) {
	switch i {
	case "pub":
		return user.Pub, nil
	case "sub":
		return user.Sub, nil
	case "pubsub":
		return user.PubSub, nil
	}
	return "", fmt.Errorf("%v is a invalid acces type", i)
}
