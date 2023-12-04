package authenticator

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidKeyType = errors.New("cannot determine the key type")
	ErrNoKeys         = errors.New("at least one key required")
)

func (b Builder) GenerateKeys(method string, keys map[string]string) (map[string]any, error) {
	var (
		keyList map[string]any
		err     error
	)

	// https://jwt.io/
	switch {
	case strings.HasPrefix(method, "RS"):
		keyList, err = b.GenerateRSAKeys(keys)
	case strings.HasPrefix(method, "HS"):
		keyList, err = b.GenerateHMACKeys(keys)
	case strings.HasPrefix(method, "ES"):
		keyList, err = b.GenerateECDSAKeys(keys)
	default:
		return nil, ErrInvalidKeyType
	}

	if err != nil {
		return nil, fmt.Errorf("reading keys failed %w", err)
	}

	if len(keyList) == 0 {
		return nil, ErrNoKeys
	}

	return keyList, nil
}

func (b Builder) GenerateRSAKeys(raw map[string]string) (map[string]any, error) {
	keys := make(map[string]any)

	for iss, publicKey := range raw {
		bytes, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			return nil, fmt.Errorf("could not read public key %w", err)
		}

		keys[iss] = bytes
	}

	return keys, nil
}

func (b Builder) GenerateECDSAKeys(raw map[string]string) (map[string]any, error) {
	keys := make(map[string]any)

	for iss, publicKey := range raw {
		bytes, err := jwt.ParseECPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			return nil, fmt.Errorf("could not read public key %w", err)
		}

		keys[iss] = bytes
	}

	return keys, nil
}

func (b Builder) GenerateHMACKeys(raw map[string]string) (map[string]any, error) {
	keys := make(map[string]any)

	for iss, key := range raw {
		bytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return nil, fmt.Errorf("failed to generate hmac key from base64 %w", err)
		}

		keys[iss] = bytes
	}

	return keys, nil
}
