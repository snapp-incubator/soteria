package config

import (
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/logger"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/tracing"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
)

const (
	DefaultHTTPPort            = 9999
	DefaultDriverHashLength    = 15
	DefaultPassengerHashLength = 15
)

// Default return default configuration.
// nolint: funlen, gomnd
func Default() Config {
	return Config{
		AllowedAccessTypes: []string{
			"pub",
			"sub",
		},
		PassengerHashLength: DefaultPassengerHashLength,
		DriverHashLength:    DefaultDriverHashLength,
		PassengerSalt:       "secret",
		DriverSalt:          "secret",
		JWT: &JWT{
			Path: "test/",
		},
		Logger: logger.Config{
			Level: "warn",
		},
		HTTPPort: DefaultHTTPPort,
		Tracer: tracing.Config{
			Enabled: false,
			Ratio:   0.1,
			Agent: tracing.Agent{
				Host: "127.0.0.1",
				Port: "6831",
			},
		},
		Company: "snapp",
		Topics: []topics.Topic{
			{
				Type:     "cab_event",
				Template: "{{.audience}}-event-{{.hashId}}",
				HashType: topics.MD5,
				Accesses: map[string]acl.AccessType{
					topics.Driver:    acl.Sub,
					topics.Passenger: acl.Sub,
				},
			},
			{
				Type:     "driver_location",
				Template: "{{.company}}/driver/{{.hashId}}/location",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.Driver:    acl.Pub,
					topics.Passenger: acl.None,
				},
			},
			{
				Type:     "passenger_location",
				Template: "{{.company}}/passenger/{{.hashId}}/location",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.Driver:    acl.Pub,
					topics.Passenger: acl.Pub,
				},
			},
			{
				Type:     "superapp_event",
				Template: "{{.company}}/{{.audience}}/{{.hashId}}/superapp",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.Driver:    acl.Sub,
					topics.Passenger: acl.Sub,
				},
			},
			{
				Type:     "box_event",
				Template: "bucks",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.Driver:    acl.None,
					topics.Passenger: acl.None,
				},
			},
			{
				Type:     "shared_location",
				Template: "{{.company}}/{{.audience}}/{{.hashId}}/{{.peer}}-location",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.Driver:    acl.Sub,
					topics.Passenger: acl.Sub,
				},
			},
			{
				Type:     "chat",
				Template: "{{.company}}/{{.audience}}/{{.hashId}}/chat",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.Driver:    acl.Sub,
					topics.Passenger: acl.Sub,
				},
			},
			{
				Type:     "general_call_entry",
				Template: "shared/{{.company}}/{{.audience}}/{{.hashId}}/call/send",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.Driver:    acl.Pub,
					topics.Passenger: acl.Pub,
				},
			},
			{
				Type:     "node_call_entry",
				Template: "{{.company}}/{{.audience}}/{{.hashId}}/call/[a-zA-Z0-9-]+/send",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.Driver:    acl.Pub,
					topics.Passenger: acl.Pub,
				},
			},
			{
				Type:     "call_outgoing",
				Template: "{{.company}}/{{.audience}}/{{.hashId}}/call/receive",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.Driver:    acl.Sub,
					topics.Passenger: acl.Sub,
				},
			},
		},
	}
}
