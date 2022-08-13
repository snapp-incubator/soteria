package authenticator

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	"go.uber.org/zap"
)

type Builder struct {
	Vendors []config.Vendor
	Logger  zap.Logger
}

func (b Builder) Authenticators() map[string]*Authenticator {
	all := make(map[string]*Authenticator)

	for _, vendor := range b.Vendors {
		publicKeys := b.PublicKeys(vendor.Keys)
		hid := HIDManager(vendor.DriverSalt, vendor.DriverHashLength, vendor.PassengerSalt, vendor.PassengerHashLength)
		allowedAccessTypes := b.GetAllowedAccessTypes(vendor.AllowedAccessTypes)

		auth := &Authenticator{
			PublicKeys:         publicKeys,
			AllowedAccessTypes: allowedAccessTypes,
			Company:            vendor.Company,
			TopicManager:       topics.NewTopicManager(vendor.Topics, hid, vendor.Company, vendor.IssEntityMap),
		}

		all[vendor.Company] = auth
	}

	return all
}

func (b Builder) PublicKeys(keys map[string]string) map[string]*rsa.PublicKey {
	rsaKeys := make(map[string]*rsa.PublicKey)

	for iss, publicKey := range keys {
		rsaKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			b.Logger.Fatal("could not read public key", zap.String("issuer", iss))
		}

		rsaKeys[iss] = rsaKey
	}

	return rsaKeys
}

// GetAllowedAccessTypes will return all allowed access types in Soteria.
func (b Builder) GetAllowedAccessTypes(accessTypes []string) []acl.AccessType {
	allowedAccessTypes := make([]acl.AccessType, 0, len(accessTypes))

	for _, a := range accessTypes {
		at, err := toUserAccessType(a)
		if err != nil {
			err = fmt.Errorf("could not convert %s: %w", at, err)
			b.Logger.Fatal("error while getting allowed access types", zap.Error(err))
		}

		allowedAccessTypes = append(allowedAccessTypes, at)
	}

	return allowedAccessTypes
}

func HIDManager(
	driverSalt string,
	driverHashLength int,
	passengerSalt string,
	passengerHashLength int,
) *snappids.HashIDSManager {
	return &snappids.HashIDSManager{
		Salts: map[snappids.Audience]string{
			snappids.DriverAudience:    driverSalt,
			snappids.PassengerAudience: passengerSalt,
		},
		Lengths: map[snappids.Audience]int{
			snappids.DriverAudience:    driverHashLength,
			snappids.PassengerAudience: passengerHashLength,
		},
	}
}

// toUserAccessType will convert string access type to it's own type.
// nolint: goerr113
func toUserAccessType(access string) (acl.AccessType, error) {
	switch access {
	case "pub":
		return acl.Pub, nil
	case "sub":
		return acl.Sub, nil
	case "pubsub":
		return acl.PubSub, nil
	}

	return "", fmt.Errorf("%v is a invalid acces type", access)
}
