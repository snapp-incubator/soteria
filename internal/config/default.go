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
		Topics: map[string]string{
			"cab_event":          `(\w+)-event-[a-zA-Z0-9]+`,
			"driver_location":    `/driver/[a-zA-Z0-9]+/location`,
			"passenger_location": `/passenger/[a-zA-Z0-9]+/location`,
			"superapp_event":     `/(driver|passenger)/[a-zA-Z0-9]+/(superapp)`,
			"box_event":          `bucks`,
			"shared_location":    `/(driver|passenger)/[a-zA-Z0-9]+/(driver-location|passenger-location)`,
			"chat":               `/(driver|passenger)/[a-zA-Z0-9]+/chat`,
			"general_call_entry": `/(driver|passenger)/[a-zA-Z0-9]+/call/send`,
			"node_call_entry":    `/(driver|passenger)/[a-zA-Z0-9]+/call/[a-zA-Z0-9-]+/send`,
			"call_outgoing":      `/(driver|passenger)/[a-zA-Z0-9]+/call/receive`,
		},
	}
}
