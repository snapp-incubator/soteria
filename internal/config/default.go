package config

import (
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/logger"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/tracing"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
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
		Users: []user.User{
			{
				Username: string(user.Driver),
				Rules: []user.Rule{
					{
						Topic:  topics.DriverLocation,
						Access: acl.Pub,
					},
					{
						Topic:  topics.CabEvent,
						Access: acl.Sub,
					},
					{
						Topic:  topics.SuperappEvent,
						Access: acl.Sub,
					},
					{
						Topic:  topics.PassengerLocation,
						Access: acl.Pub,
					},
					{
						Topic:  topics.SharedLocation,
						Access: acl.Sub,
					},
					{
						Topic:  topics.Chat,
						Access: acl.Sub,
					},
					{
						Topic:  topics.GeneralCallEntry,
						Access: acl.Pub,
					},
					{
						Topic:  topics.NodeCallEntry,
						Access: acl.Pub,
					},
					{
						Topic:  topics.CallOutgoing,
						Access: acl.Sub,
					},
				},
			},
			{
				Username: string(user.Passenger),
				Rules: []user.Rule{
					{
						Topic:  topics.CabEvent,
						Access: acl.Sub,
					},
					{
						Topic:  topics.SuperappEvent,
						Access: acl.Sub,
					},
					{
						Topic:  topics.PassengerLocation,
						Access: acl.Pub,
					},
					{
						Topic:  topics.SharedLocation,
						Access: acl.Sub,
					},
					{
						Topic:  topics.Chat,
						Access: acl.Sub,
					},
					{
						Topic:  topics.GeneralCallEntry,
						Access: acl.Pub,
					},
					{
						Topic:  topics.NodeCallEntry,
						Access: acl.Pub,
					},
					{
						Topic:  topics.CallOutgoing,
						Access: acl.Sub,
					},
				},
			},
		},
		Topics: []topics.Topic{
			{
				Type:     "cab_event",
				Template: "{{.audience}}-event-{{.hashId}}",
				Regex:    `(\w+)-event-[a-zA-Z0-9]+`,
			},
			{
				Type:     "driver_location",
				Template: "{{.company}}/driver/{{.hashId}}/location",
				Regex:    `/driver/[a-zA-Z0-9]+/location`,
			},
			{
				Type:     "passenger_location",
				Template: "{{.company}}/passenger/{{.hashId}}/location",
				Regex:    "/passenger/[a-zA-Z0-9]+/location",
			},
			{
				Type:     "superapp_event",
				Template: "{{.company}}/{{.audience}}/{{.hashId}}/superapp",
				Regex:    `/(driver|passenger)/[a-zA-Z0-9]+/(superapp)`,
			},
			{
				Type:     "box_event",
				Template: "bucks",
				Regex:    "bucks",
			},
			{
				Type:     "shared_location",
				Template: "{{.company}}/{{.audience}}/{{.hashId}}/{{.peer}}-location",
				Regex:    `/(driver|passenger)/[a-zA-Z0-9]+/(driver-location|passenger-location)`,
			},
			{
				Type:     "chat",
				Template: "{{.company}}/{{.audience}}/{{.hashId}}/chat",
				Regex:    `/(driver|passenger)/[a-zA-Z0-9]+/chat`,
			},
			{
				Type:     "general_call_entry",
				Template: "shared/{{.company}}/{{.audience}}/{{.hashId}}/call/send",
				Regex:    `/(driver|passenger)/[a-zA-Z0-9]+/call/send`,
			},
			{
				Type:     "node_call_entry",
				Template: "{{.company}}/{{.audience}}/{{.hashId}}/call/{{.node}}/send",
				Regex:    `/(driver|passenger)/[a-zA-Z0-9]+/call/[a-zA-Z0-9-]+/send`,
			},
			{
				Type:     "call_outgoing",
				Template: "{{.company}}/{{.audience}}/{{.hashId}}/call/receive",
				Regex:    `/(driver|passenger)/[a-zA-Z0-9]+/call/receive`,
			},
		},
	}
}
