package config

import "gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"

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
				Username: "driver",
				Rules:    []user.Rule{},
			},
			{
				Username: "passenger",
				Rules:    []user.Rule{},
			},
		},
	}
}
