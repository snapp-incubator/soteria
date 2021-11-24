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
// nolint: funlen
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
						Topic:  topics.CallEntry,
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
						Topic:  topics.CallEntry,
						Access: acl.Pub,
					},
					{
						Topic:  topics.CallOutgoing,
						Access: acl.Sub,
					},
				},
			},
		},
	}
}
