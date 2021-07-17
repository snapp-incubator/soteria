package acl

import "net"

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

// ValidateIP takes validIPs and invalidIPs and tell whether a IP is valid or not.
func ValidateIP(ip string, validIPs, invalidIPs []string) bool {
	isValid := false

	for _, validIP := range validIPs {
		_, network, err := net.ParseCIDR(validIP)
		if err != nil && validIP == ip {
			isValid = true

			break
		} else if err == nil && network.Contains(net.ParseIP(ip)) {
			isValid = true

			break
		}
	}

	for _, invalidIP := range invalidIPs {
		_, network, err := net.ParseCIDR(invalidIP)
		if err != nil && invalidIP == ip {
			isValid = false

			break
		} else if err == nil && network.Contains(net.ParseIP(ip)) {
			isValid = false

			break
		}
	}

	return isValid
}
