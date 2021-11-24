package config

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/logger"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
)

const (
	// Prefix indicates environment variables prefix.
	Prefix = "soteria_"
)

type (
	// Config is the main container of Soteria's config.
	Config struct {
		AllowedAccessTypes  []string      `koanf:"allowed_access_types"`
		PassengerHashLength int           `koanf:"passenger_hash_length"`
		DriverHashLength    int           `koanf:"driver_hash_length"`
		PassengerSalt       string        `koanf:"passenger_salt"`
		DriverSalt          string        `koanf:"driver_salt"`
		JWT                 *JWT          `koanf:"jwt"`
		Logger              logger.Config `koanf:"logger"`
		HTTPPort            int           `koanf:"http_port"`
		Tracer              *TracerConfig `koanf:"tracer"`
		Company             string        `koanf:"company"`
		Users               []user.User   `koanf:"users"`
	}

	// JWt contains path of the keys for JWT encryption.
	JWT struct {
		Path string `koanf:"path"`
	}

	// Tracer contains all configs needed to create a tracer.
	TracerConfig struct {
		Enabled      bool    `koanf:"enabled"`
		ServiceName  string  `koanf:"service_name"`
		SamplerType  string  `koanf:"sampler_type"`
		SamplerParam float64 `koanf:"sampler_param"`
		Host         string  `koanf:"host"`
		Port         int     `koanf:"port"`
	}
)

// New reads configuration with koanf.
func New() Config {
	var instance Config

	k := koanf.New(".")

	// load default configuration from file
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load configuration from file
	if err := k.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		log.Printf("error loading config.yml: %s", err)
	}

	// load environment variables
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, Prefix)), "_", ".")
	}), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	log.Printf("following configuration is loaded:\n%+v", instance)

	return instance
}

// ReadPrivateKey will read and return private key that is used for JWT encryption.
func (a *Config) ReadPublicKey(u user.Issuer) (*rsa.PublicKey, error) {
	var fileName string

	switch u {
	case user.Driver:
		fileName = fmt.Sprintf("%s%s", a.JWT.Path, "0.pem")
	case user.Passenger:
		fileName = fmt.Sprintf("%s%s", a.JWT.Path, "1.pem")
	default:
		return nil, errors.New("invalid issuer, public key not found")
	}

	pem, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pem)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

// GetAllowedAccessTypes will return all allowed access types in Soteria.
func (a *Config) GetAllowedAccessTypes() ([]acl.AccessType, error) {
	allowedAccessTypes := make([]acl.AccessType, 0, len(a.AllowedAccessTypes))

	for _, a := range a.AllowedAccessTypes {
		at, err := toUserAccessType(a)
		if err != nil {
			return nil, fmt.Errorf("could not convert %s: %w", at, err)
		}

		allowedAccessTypes = append(allowedAccessTypes, at)
	}

	return allowedAccessTypes, nil
}

// toUserAccessType will convert string access type to it's own type.
func toUserAccessType(access string) (acl.AccessType, error) {
	switch access {
	case "pub":
		return acl.Pub, nil
	case "sub":
		return acl.Sub, nil
	case "pubsub":
		return acl.PubSub, nil
	}

	return "", fmt.Errorf("%v is a invalid acces type", access)
}
