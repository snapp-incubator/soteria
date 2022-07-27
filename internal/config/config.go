package config

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/tidwall/pretty"
	"gitlab.snapp.ir/dispatching/soteria/internal/logger"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/internal/tracing"
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

	Vendor struct {
		AllowedAccessTypes  []string       `koanf:"allowed_access_types"`
		PassengerHashLength int            `koanf:"passenger_hash_length"`
		DriverHashLength    int            `koanf:"driver_hash_length"`
		PassengerSalt       string         `koanf:"passenger_salt"`
		DriverSalt          string         `koanf:"driver_salt"`
		Company             string         `koanf:"company"`
		Topics              []topics.Topic `koanf:"topics"`
		DriverKey           string         `koanf:"driver_key"`
		PassengerKey        string         `koanf:"passenger_key"`
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
