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
		HashIDMap: map[string]topics.HashData{
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
				Template: "^{{IssToEntity .iss}}-event-{{ EncodeMD5 (DecodeHashID .sub .iss) }}$",
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Sub,
					topics.PassengerIss: acl.Sub,
				},
			},
			{
				Type:     topics.DriverLocation,
				Template: "^{{.company}}/driver/{{.sub}}/location$",
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Pub,
					topics.PassengerIss: acl.None,
				},
			},
			{
				Type:     topics.PassengerLocation,
				Template: "^{{.company}}/passenger/{{.sub}}/location$",
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Pub,
					topics.PassengerIss: acl.Pub,
				},
			},
			{
				Type:     topics.SuperappEvent,
				Template: "^{{.company}}/{{IssToEntity .iss}}/{{.sub}}/superapp$",
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Sub,
					topics.PassengerIss: acl.Sub,
				},
			},
			{
				Type:     topics.BoxEvent,
				Template: "^bucks$",
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.None,
					topics.PassengerIss: acl.None,
				},
			},
			{
				Type:     topics.SharedLocation,
				Template: "^{{.company}}/{{IssToEntity .iss}}/{{.sub}}/{{IssToPeer .iss}}-location$", //nolint:lll
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Sub,
					topics.PassengerIss: acl.Sub,
				},
			},
			{
				Type:     topics.Chat,
				Template: "^{{.company}}/{{IssToEntity .iss}}/{{.sub}}/chat$",
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Sub,
					topics.PassengerIss: acl.Sub,
				},
			},
			{
				Type:     topics.GeneralCallEntry,
				Template: "^shared/{{.company}}/{{IssToEntity .iss}}/{{.sub}}/call/send$",
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Pub,
					topics.PassengerIss: acl.Pub,
				},
			},
			{
				Type:     topics.NodeCallEntry,
				Template: "^{{.company}}/{{IssToEntity .iss}}/{{.sub}}/call/[a-zA-Z0-9-_]+/send$", //nolint: lll
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Pub,
					topics.PassengerIss: acl.Pub,
				},
			},
			{
				Type:     topics.CallOutgoing,
				Template: "^{{.company}}/{{IssToEntity .iss}}/{{.sub}}/call/receive$",
				Accesses: map[string]acl.AccessType{
					topics.DriverIss:    acl.Sub,
					topics.PassengerIss: acl.Sub,
				},
			},
		},
		// 		Keys: map[string]string{
		// 			"0": `-----BEGIN PUBLIC KEY-----
		// MIIBITANBgkqhkiG9w0BAQEFAAOCAQ4AMIIBCQKCAQBk7O6M5p4eYNAwtVU2beGa
		// W4mhFG94OtYUWDl1E7UUrhUNGf97Eb/45NjQszu0YPERnApJc2RUm2TrS7iq0mHz
		// Xbwf+CbNF54Q5mjuHcpBKgvFwUUSCCYBftmRc4xbFIH4Oh3nHC2GeukUS9TmJwjM
		// tJKyU0Ve8BK5BgjhagM7XSs+scE2mxemoWtcs6mJLtBuEgRGMgHW00mSdOcLp/+l
		// oHpSzRYN92/DomwmmjGVy8Ji0faeHx+r79ZzE0E8Rcc29Yhrg1ymrjfkXg98WjAb
		// TSv4UAN20lsBDejpnGEZKJrxHZ56gHgaJn6PKKCD6ItJA7y7iraCdBhCfAIUIz/z
		// AgMBAAE=
		// -----END PUBLIC KEY-----`,

		// 			"1": `-----BEGIN PUBLIC KEY-----
		// MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1lNRwyNsDieWs6LvHOJ+
		// GyehhRC4Pn5yL5edKP3565F3LtRDMrkzwDRsQbqnUtTea9HCdTdBv+lI8vE17qRi
		// RQn10IMaIH6e4Aa3OWNClFhuqNOag7VmffsjTOgxHgHpfGAKVF/4BwqOHrdHFbAD
		// VOiWB1hv9Uc0C5laffGAub7fj+EAI02zlrsNDxYW8vyF2H47N7VWcvgd3RhZpxlG
		// 8bq9phl7Ja55YmQiT2Ic3/K5tsazg5z9lz6OTrx+JvWbefHFlJpjCLz5yefEaRmX
		// 9L/zyDMi4jgFTZEWNXC2vIrxwZMFwFhBXEp0PcCbuHJgJIucbRrbwukQC16uHJwP
		// zQIDAQAB
		// -----END PUBLIC KEY-----`,
		// 		},
		Keys: map[string][]string{
			"0": []string{`-----BEGIN PUBLIC KEY-----
	MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAuiAp6wih+q53HJQhcQed
	1ergtAP+keqrtRdtUMo3wPEupVZILd+JCMs3MV9yxLxLD55J9pESNhqiJjeXHaCi
	w8d0y5aKmoa3cjHqG2vzKH7R8n8lWPy++oGvEhyYVtQyQmOOMgZ2jyRAdk5Co5iK
	2f5O0ydnTobenbo2u19d0DGJVULynYEsuOQJEHs+z+l7DoiaTba6hrxAWWncUVWd
	emplkhWX7oyujgfIClaSr15hWxZK8jSPrsJE0vbMmWKGu3LbW3aVRfWKwnRVdJLM
	zhpLNEZiK1CnG/OKhXqmT6/n/MuqTzf168I1J9ypPsWpLiU/jc7C1weeut3LhFTt
	rwUrADTMRUD96F6pk5gIjDm9gvSPmkGRYkq+6og2MhYFAXNwF1t8c9Ht8V93JlrF
	zmLNfQ1yTAvwVz1ba7PwzqkJk2U0nCMCQMfNxhCeS6uXBdtqQXHeevht7frOkSqa
	tF/jU2wcesAX2XUv/Hg9X+eYvK8KN7iN2sQ+4WuwheqOuIsazTSAk93+YEZx+UkT
	00AlIUi8IbmHReBAhxOTPzodFVR7jzyLfNRB0n5dbStAYhjK2QxgbNF6tRcfyT25
	kYfTXFJg8TGMIjUJpA0/JVRSnewpSm2jFh8s2RaaAy27IAnHmN8G5tJrEJ05Vcp6
	KNke4qsRTcXvt397Z9a6fDECAwEAAQ==
-----END PUBLIC KEY-----`, `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAuiAp6wih+q53HJQhcQed
1ergtAP+keqrtRdtUMo3wPEupVZILd+JCMs3MV9yxLxLD55J9pESNhqiJjeXHaCi
w8d0y5aKmoa3cjHqG2vzKH7R8n8lWPy++oGvEhyYVtQyQmOOMgZ2jyRAdk5Co5iK
2f5O0ydnTobenbo2u19d0DGJVULynYEsuOQJEHs+z+l7DoiaTba6hrxAWWncUVWd
emplkhWX7oyujgfIClaSr15hWxZK8jSPrsJE0vbMmWKGu3LbW3aVRfWKwnRVdJLM
zhpLNEZiK1CnG/OKhXqmT6/n/MuqTzf168I1J9ypPsWpLiU/jc7C1weeut3LhFTt
rwUrADTMRUD96F6pk5gIjDm9gvSPmkGRYkq+6og2MhYFAXNwF1t8c9Ht8V93JlrF
zmLNfQ1yTAvwVz1ba7PwzqkJk2U0nCMCQMfNxhCeS6uXBdtqQXHeevht7frOkSqa
tF/jU2wcesAX2XUv/Hg9X+eYvK8KN7iN2sQ+4WuwheqOuIsazTSAk93+YEZx+UkT
00AlIUi8IbmHReBAhxOTPzodFVR7jzyLfNRB0n5dbStAYhjK2QxgbNF6tRcfyT25
kYfTXFJg8TGMIjUJpA0/JVRSnewpSm2jFh8s2RaaAy27IAnHmN8G5tJrEJ05Vcp6
KNke4qsRTcXvt397Z9a6fDECAwEAAQ==
-----END PUBLIC KEY-----`},

			"1": []string{`-----BEGIN PUBLIC KEY-----
	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCuxK5lFw6U/aHSSbxaCEsfxqLA
	zqQeioXV8QxbectY95p7OoCPBMyxZ6FXipbzmNLEO41kjFCHthAZ4DHhx6q/bEF3
	0sj9J2FwL3rO3mc31hbyUAGaIjTgR4302MXgnTmeX68dOqmgBZem70Si8gvbgXoc
	qF+zAHZiEZ4hr24/KQIDAQAB
	-----END PUBLIC KEY-----`, `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCuxK5lFw6U/aHSSbxaCEsfxqLA
zqQeioXV8QxbectY95p7OoCPBMyxZ6FXipbzmNLEO41kjFCHthAZ4DHhx6q/bEF3
0sj9J2FwL3rO3mc31hbyUAGaIjTgR4302MXgnTmeX68dOqmgBZem70Si8gvbgXoc
qF+zAHZiEZ4hr24/KQIDAQAB
-----END PUBLIC KEY-----`},
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
