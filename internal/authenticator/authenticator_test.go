package authenticator_test

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// nolint: gosec, lll
	invalidToken = "ey1JhbGciOiJSUzI1NiIsInR5cCI56kpXVCJ9.eyJzdWIiOiJCRzdScDFkcnpWRE5RcjYiLCJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwiaXNzIjowLCJpYXQiOjE1MTYyMzkwMjJ9.1cYXFEhcewOYFjGJYhB8dsaFO9uKEXwlM8954rkt4Tsu0lWMITbRf_hHh1l9QD4MFqD-0LwRPUYaiaemy0OClMu00G2sujLCWaquYDEP37iIt8RoOQAh8Jb5vT8LX5C3PEKvbW_i98u8HHJoFUR9CXJmzrKi48sAcOYvXVYamN0S9KoY38H-Ze37Mdu3o6B58i73krk7QHecsc2_PkCJisvUVAzb0tiInIalBc8-zI3QZSxwNLr_hjlBg1sUxTUvH5SCcRR7hxI8TxJzkOHqAHWDRO84NC_DSAoO2p04vrHpqglN9XPJ8RC2YWpfefvD2ttH554RJWu_0RlR2kAYvQ"

	validPassengerCabEventTopic   = "passenger-event-152384980615c2bd16143cff29038b67"
	invalidPassengerCabEventTopic = "passenger-event-152384980615c2bd16156cff29038b67"

	validDriverCabEventTopic   = "driver-event-152384980615c2bd16143cff29038b67"
	invalidDriverCabEventTopic = "driver-event-152384980615c2bd16156cff29038b67"

	validDriverLocationTopic   = "snapp/driver/DXKgaNQa7N5Y7bo/location"
	invalidDriverLocationTopic = "snapp/driver/DXKgaNQa9Q5Y7bo/location"

	validPassengerSuperappEventTopic   = "snapp/passenger/DXKgaNQa7N5Y7bo/superapp"
	invalidPassengerSuperappEventTopic = "snapp/passenger/DXKgaNQa9Q5Y7bo/superapp"

	validDriverSuperappEventTopic   = "snapp/driver/DXKgaNQa7N5Y7bo/superapp"
	invalidDriverSuperappEventTopic = "snapp/driver/DXKgaNQa9Q5Y7bo/superapp"

	validDriverSharedTopic      = "snapp/driver/DXKgaNQa7N5Y7bo/passenger-location"
	validPassengerSharedTopic   = "snapp/passenger/DXKgaNQa7N5Y7bo/driver-location"
	invalidDriverSharedTopic    = "snapp/driver/0596923be632d673560af9adadd2f78a/passenger-location"
	invalidPassengerSharedTopic = "snapp/passenger/0596923be632d673560af9adadd2f78a/driver-location"

	validDriverChatTopic      = "snapp/driver/DXKgaNQa7N5Y7bo/chat"
	validPassengerChatTopic   = "snapp/passenger/DXKgaNQa7N5Y7bo/chat"
	invalidDriverChatTopic    = "snapp/driver/0596923be632d673560af9adadd2f78a/chat"
	invalidPassengerChatTopic = "snapp/passenger/0596923be632d673560af9adadd2f78a/chat"

	validDriverCallEntryTopic         = "shared/snapp/driver/DXKgaNQa7N5Y7bo/call/send"
	validPassengerCallEntryTopic      = "shared/snapp/passenger/DXKgaNQa7N5Y7bo/call/send"
	validDriverNodeCallEntryTopic     = "snapp/driver/DXKgaNQa7N5Y7bo/call/heliograph-0/send"
	validPassengerNodeCallEntryTopic  = "snapp/passenger/DXKgaNQa7N5Y7bo/call/heliograph-0/send"
	invalidDriverCallEntryTopic       = "snapp/driver/0596923be632d673560af9adadd2f78a/call/send"
	invalidPassengerCallEntryTopic    = "snapp/passenger/0596923be632d673560af9adadd2f78a/call/send"
	validDriverCallOutgoingTopic      = "snapp/driver/DXKgaNQa7N5Y7bo/call/receive"
	validPassengerCallOutgoingTopic   = "snapp/passenger/DXKgaNQa7N5Y7bo/call/receive"
	invalidDriverCallOutgoingTopic    = "snapp/driver/0596923be632d673560af9adadd2f78a/call/receive"
	invalidPassengerCallOutgoingTopic = "snapp/passenger/0596923be632d673560af9adadd2f78a/call/receive"
)

var (
	ErrPrivateKeyNotFound = errors.New("invalid user, private key not found")
	ErrPublicKeyNotFound  = errors.New("invalid user, public key not found")
)

func getPublicKey(u string) (*rsa.PublicKey, error) {
	var fileName string

	switch u {
	case "1":
		fileName = "../../test/snapp-1.pem"
	case "0":
		fileName = "../../test/snapp-0.pem"
	case "admin":
		fileName = "../../test/snapp-admin.pem"
	default:
		return nil, ErrPublicKeyNotFound
	}

	pem, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("reading public key failed %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pem)
	if err != nil {
		return nil, fmt.Errorf("paring public key failed %w", err)
	}

	return publicKey, nil
}

func getPrivateKey(u string) (*rsa.PrivateKey, error) {
	var fileName string

	switch u {
	case "0":
		fileName = "../../test/snapp-0.private.pem"
	case "1":
		fileName = "../../test/snapp-1.private.pem"
	case "admin":
		fileName = "../../test/snapp-admin.private.pem"
	default:
		return nil, ErrPrivateKeyNotFound
	}

	pem, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("reading private key failed %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return nil, fmt.Errorf("paring private key failed %w", err)
	}

	return privateKey, nil
}

func getSampleToken(issuer string, key *rsa.PrivateKey) (string, error) {
	exp := time.Now().Add(time.Hour * 24 * 365 * 10)
	sub := "DXKgaNQa7N5Y7bo"

	// nolint: exhaustruct
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(exp),
		Issuer:    issuer,
		Subject:   sub,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("cannot generate a signed string %w", err)
	}

	return tokenString, nil
}
