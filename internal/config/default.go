package config

import (
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
)

// Default return default configuration.
func Default() Config {
	return Config{
		AllowedAccessTypes: []string{
			"pub",
			"sub",
		},
		PassengerHashLength: 15,
		DriverHashLength:    15,
		PassengerSalt:       "secret",
		DriverSalt:          "secret",
		JWT: &JWT{
			Path: "/test",
		},
		Logger: &Logger{
			Level: "warn",
		},
		HTTPPort: 0,
		Tracer: &TracerConfig{
			Enabled:      false,
			ServiceName:  "",
			SamplerType:  "",
			SamplerParam: 0.0,
			Host:         "",
			Port:         0,
		},
		Company: "snapp",
		Users: []user.User{
			{
				Username: string(user.Driver),
				Rules: []user.Rule{
					{
						Topic:      topics.DriverLocation,
						AccessType: acl.Pub,
					},
					{
						Topic:      topics.CabEvent,
						AccessType: acl.Sub,
					},
					{
						Topic:      topics.SuperappEvent,
						AccessType: acl.Sub,
					},
					{
						Topic:      topics.PassengerLocation,
						AccessType: acl.Pub,
					},
					{
						Topic:      topics.SharedLocation,
						AccessType: acl.Sub,
					},
					{
						Topic:      topics.Chat,
						AccessType: acl.Sub,
					},
					{
						Topic:      topics.CallEntry,
						AccessType: acl.Pub,
					},
					{
						Topic:      topics.CallOutgoing,
						AccessType: acl.Sub,
					},
				},
			},
			{
				Username: string(user.Passenger),
				Rules: []user.Rule{
					{
						Topic:      topics.CabEvent,
						AccessType: acl.Sub,
					},
					{
						Topic:      topics.SuperappEvent,
						AccessType: acl.Sub,
					},
					{
						Topic:      topics.PassengerLocation,
						AccessType: acl.Pub,
					},
					{
						Topic:      topics.SharedLocation,
						AccessType: acl.Sub,
					},
					{
						Topic:      topics.Chat,
						AccessType: acl.Sub,
					},
					{
						Topic:      topics.CallEntry,
						AccessType: acl.Pub,
					},
					{
						Topic:      topics.CallOutgoing,
						AccessType: acl.Sub,
					},
				},
			},
		},
	}
}
