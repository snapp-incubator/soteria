package config

import (
	"gitlab.snapp.ir/dispatching/soteria/internal/logger"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/internal/tracing"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
)

const (
	DefaultHTTPPort            = 9999
	DefaultDriverHashLength    = 15
	DefaultPassengerHashLength = 15
)

// Default return default configuration.
// nolint: gomnd
func Default() Config {
	return Config{
		DefaultVendor: "snapp",
		Vendors: []Vendor{
			SnappVendor(),
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
	}
}

//nolint:funlen
func SnappVendor() Vendor {
	return Vendor{
		AllowedAccessTypes: []string{
			"pub",
			"sub",
		},
		HashIDMap: map[string]HashID{
			"0": {
				Alphabet: "",
				Length:   DefaultDriverHashLength,
				Salt:     "secret",
			},
			"1": {
				Alphabet: "",
				Length:   DefaultPassengerHashLength,
				Salt:     "secret",
			},
		},
		Company: "snapp",
		Topics: []topics.Topic{
			{
				Type:     topics.CabEvent,
				Template: "^{{IssToEntity .iss}}-event-{{HashID .hashType .sub .iss}}$",
				HashType: topics.MD5,
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Sub,
					topics.PassengerIss: acl.Sub,
				},
			},
			{
				Type:     topics.DriverLocation,
				Template: "^{{.company}}/driver/{{HashID .hashType .sub .iss}}/location$",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Pub,
					topics.PassengerIss: acl.None,
				},
			},
			{
				Type:     topics.PassengerLocation,
				Template: "^{{.company}}/passenger/{{HashID .hashType .sub .iss}}/location$",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Pub,
					topics.PassengerIss: acl.Pub,
				},
			},
			{
				Type:     topics.SuperappEvent,
				Template: "^{{.company}}/{{IssToEntity .iss}}/{{HashID .hashType .sub .iss}}/superapp$",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Sub,
					topics.PassengerIss: acl.Sub,
				},
			},
			{
				Type:     topics.BoxEvent,
				Template: "^bucks$",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.None,
					topics.PassengerIss: acl.None,
				},
			},
			{
				Type:     topics.SharedLocation,
				Template: "^{{.company}}/{{IssToEntity .iss}}/{{HashID .hashType .sub .iss}}/{{IssToPeer .iss}}-location$", //nolint:lll
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Sub,
					topics.PassengerIss: acl.Sub,
				},
			},
			{
				Type:     topics.Chat,
				Template: "^{{.company}}/{{IssToEntity .iss}}/{{HashID .hashType .sub .iss}}/chat$",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Sub,
					topics.PassengerIss: acl.Sub,
				},
			},
			{
				Type:     topics.GeneralCallEntry,
				Template: "^shared/{{.company}}/{{IssToEntity .iss}}/{{HashID .hashType .sub .iss}}/call/send$",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Pub,
					topics.PassengerIss: acl.Pub,
				},
			},
			{
				Type:     topics.NodeCallEntry,
				Template: "^{{.company}}/{{IssToEntity .iss}}/{{HashID .hashType .sub .iss}}/call/[a-zA-Z0-9-_]+/send$", //nolint: lll
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Pub,
					topics.PassengerIss: acl.Pub,
				},
			},
			{
				Type:     topics.CallOutgoing,
				Template: "^{{.company}}/{{IssToEntity .iss}}/{{HashID .hashType .sub .iss}}/call/receive$",
				HashType: topics.HashID,
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Sub,
					topics.PassengerIss: acl.Sub,
				},
			},
		},
		Keys: map[string]string{
			"0": `-----BEGIN PUBLIC KEY-----
MIIBITANBgkqhkiG9w0BAQEFAAOCAQ4AMIIBCQKCAQBk7O6M5p4eYNAwtVU2beGa
W4mhFG94OtYUWDl1E7UUrhUNGf97Eb/45NjQszu0YPERnApJc2RUm2TrS7iq0mHz
Xbwf+CbNF54Q5mjuHcpBKgvFwUUSCCYBftmRc4xbFIH4Oh3nHC2GeukUS9TmJwjM
tJKyU0Ve8BK5BgjhagM7XSs+scE2mxemoWtcs6mJLtBuEgRGMgHW00mSdOcLp/+l
oHpSzRYN92/DomwmmjGVy8Ji0faeHx+r79ZzE0E8Rcc29Yhrg1ymrjfkXg98WjAb
TSv4UAN20lsBDejpnGEZKJrxHZ56gHgaJn6PKKCD6ItJA7y7iraCdBhCfAIUIz/z
AgMBAAE=
-----END PUBLIC KEY-----`,

			"1": `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1lNRwyNsDieWs6LvHOJ+
GyehhRC4Pn5yL5edKP3565F3LtRDMrkzwDRsQbqnUtTea9HCdTdBv+lI8vE17qRi
RQn10IMaIH6e4Aa3OWNClFhuqNOag7VmffsjTOgxHgHpfGAKVF/4BwqOHrdHFbAD
VOiWB1hv9Uc0C5laffGAub7fj+EAI02zlrsNDxYW8vyF2H47N7VWcvgd3RhZpxlG
8bq9phl7Ja55YmQiT2Ic3/K5tsazg5z9lz6OTrx+JvWbefHFlJpjCLz5yefEaRmX
9L/zyDMi4jgFTZEWNXC2vIrxwZMFwFhBXEp0PcCbuHJgJIucbRrbwukQC16uHJwP
zQIDAQAB
-----END PUBLIC KEY-----`,
		},
		IssEntityMap: map[string]string{
			"0": "driver",
			"1": "passenger",
		},
		IssPeerMap: map[string]string{
			"0": "passenger",
			"1": "driver",
		},
		Jwt: Jwt{
			IssName:       "iss",
			SubName:       "sub",
			SigningMethod: "RSA512",
		},
	}
}
