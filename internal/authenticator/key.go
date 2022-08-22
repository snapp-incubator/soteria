package authenticator

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func (b Builder) GenerateKeys(method string, keys map[string]string) map[string]any {
	var keyList map[string]any

	// ES RS HS PS EdDSA
	switch {
	case strings.HasPrefix(method, "RS"):
		keyList = b.GenerateRsaKeys(keys)
	case strings.HasPrefix(method, "HS"):
		keyList = b.GenerateHMacKeys(keys)
	default:
		keyList = make(map[string]any)
	}

	return keyList
}

func (b Builder) GenerateRsaKeys(raw map[string]string) map[string]any {
	rsaKeys := make(map[string]any)

	for iss, publicKey := range raw {
		rsaKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			b.Logger.Fatal("could not read public key", zap.String("issuer", iss))
		}

		rsaKeys[iss] = rsaKey
	}

	return rsaKeys
}

func (b Builder) GenerateHMacKeys(raw map[string]string) map[string]any {
	keys := make(map[string]any)

	for iss, key := range raw {
		bytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			log.Fatalf("failed to generate hmac key: %v", err)
		}

		keys[iss] = bytes
	}

	return keys
}
