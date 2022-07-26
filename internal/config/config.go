package config

import (
	"crypto/rsa"
	"encoding/json"
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
	"github.com/tidwall/pretty"
	"gitlab.snapp.ir/dispatching/soteria/internal/logger"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/internal/tracing"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
)

const (
	// Prefix indicates environment variables prefix.
	Prefix = "soteria_"
)

type (
	// Config is the main container of Soteria's config.
	Config struct {
		Vendors  []Vendor       `koanf:"vendors"`
		Logger   logger.Config  `koanf:"logger"`
		HTTPPort int            `koanf:"http_port"`
		Tracer   tracing.Config `koanf:"tracer"`
	}

	// JWt contains path of the keys for JWT encryption.
	JWT struct {
		Path string `koanf:"path"`
	}

	Vendor struct {
		AllowedAccessTypes  []string       `koanf:"allowed_access_types"`
		PassengerHashLength int            `koanf:"passenger_hash_length"`
		DriverHashLength    int            `koanf:"driver_hash_length"`
		PassengerSalt       string         `koanf:"passenger_salt"`
		DriverSalt          string         `koanf:"driver_salt"`
		JWT                 *JWT           `koanf:"jwt"`
		Company             string         `koanf:"company"`
		Topics              []topics.Topic `koanf:"topics"`
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

	indent, err := json.MarshalIndent(instance, "", "\t")
	if err != nil {
		log.Fatalf("error marshaling configuration to json: %s", err)
	}

	indent = pretty.Color(indent, nil)
	tmpl := `
	================ Loaded Configuration ================
	%s
	======================================================
	`
	log.Printf(tmpl, string(indent))

	return instance
}

// ReadPrivateKey will read and return private key that is used for JWT encryption.
// nolint: wrapcheck, goerr113
func (a *Config) ReadPublicKey(u user.Issuer) (*rsa.PublicKey, error) {
	var fileName string

	switch u { // nolint:exhaustive
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
// nolint: goerr113
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
