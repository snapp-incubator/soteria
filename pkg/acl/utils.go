package acl

import "net"

// ValidateEndpoint takes authorizedEndpoints and unauthorizedEndpoints and tell whether a endpoint is authorized or not
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

// ValidateIP takes validIPs and invalidIPs and tell whether a IP is valid or not
func ValidateIP(IP string, validIPs, invalidIPs []string) bool {
	isValid := false
	for _, validIP := range validIPs {
		_, network, err := net.ParseCIDR(validIP)
		if err != nil && validIP == IP {
			isValid = true
			break

		} else if err == nil && network.Contains(net.ParseIP(IP)) {
			isValid = true
			break
		}
	}
	for _, invalidIP := range invalidIPs {
		_, network, err := net.ParseCIDR(invalidIP)
		if err != nil && invalidIP == IP {
			isValid = false
			break
		} else if err == nil && network.Contains(net.ParseIP(IP)) {
			isValid = false
			break
		}
	}
	return isValid
}
