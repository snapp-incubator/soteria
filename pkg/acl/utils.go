package acl

// Access Types for EMQ contains subscribe, publish and publish-subscribe.
type AccessType string

const (
	Sub    AccessType = "1"
	Pub    AccessType = "2"
	PubSub AccessType = "3"

	ClientCredentials = "client_credentials"
)

func (a AccessType) String() string {
	switch a {
	case Sub:
		return "subscribe"
	case Pub:
		return "publish"
	case PubSub:
		return "publish-subscribe"
	}

	return ""
}

// ValidateEndpoint takes authorizedEndpoints and unauthorizedEndpoints and
// tell whether a endpoint is authorized or not.
func ValidateEndpoint(endpoint string, authorizedEndpoints, unauthorizedEndpoints []string) bool {
	if len(authorizedEndpoints) == 0 && len(unauthorizedEndpoints) == 0 {
		return true
	}

	isValid := false

	for _, e := range authorizedEndpoints {
		if e == endpoint {
			isValid = true

			break
		}
	}

	for _, e := range unauthorizedEndpoints {
		if e == endpoint {
			isValid = false

			break
		}
	}

	return isValid
}
