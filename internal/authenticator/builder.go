package authenticator

import (
	"fmt"

	"github.com/speps/go-hashids/v2"
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
		b.ValidateMappers(vendor.IssEntityMap, vendor.IssPeerMap)

		keys := b.GenerateKeys(vendor.Jwt.SigningMethod, vendor.Keys)
		allowedAccessTypes := b.GetAllowedAccessTypes(vendor.AllowedAccessTypes)

		hid := make(map[string]*hashids.HashID)
		for iss, data := range vendor.HashIDMap {
			var err error

			hd := hashids.NewData()
			hd.Salt = data.Salt
			hd.MinLength = data.Length
			if data.Alphabet != "" {
				hd.Alphabet = data.Alphabet
			}

			hid[iss], err = hashids.NewWithData(hd)
			if err != nil {
				b.Logger.Fatal("cannot create hashid", zap.Error(err), zap.Any("configuration", data), zap.String("iss", iss))
			}
		}

		auth := &Authenticator{
			Keys:               keys,
			AllowedAccessTypes: allowedAccessTypes,
			Company:            vendor.Company,
			TopicManager: topics.NewTopicManager(
				vendor.Topics,
				hid,
				vendor.Company,
				vendor.IssEntityMap,
				vendor.IssPeerMap,
			),
			JwtConfig: vendor.Jwt,
		}

		all[vendor.Company] = auth
	}

	if len(all) == 0 {
		b.Logger.Fatal("at least one vendor should be enable to have soteria")
	}

	return all
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

func (b Builder) ValidateMappers(issEntityMap, issPeerMap map[string]string) {
	if _, ok := issEntityMap[topics.Default]; !ok {
		b.Logger.Fatal("default case for iss-entity map is required")
	}

	if _, ok := issPeerMap[topics.Default]; !ok {
		b.Logger.Fatal("default case for iss-peer map is required")
	}
}
