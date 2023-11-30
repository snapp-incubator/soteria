package authenticator

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/internal/topics"
	"github.com/snapp-incubator/soteria/pkg/acl"
	"github.com/snapp-incubator/soteria/pkg/validator"
	"go.uber.org/zap"
)

type Builder struct {
	Vendors         []config.Vendor
	Logger          zap.Logger
	ValidatorConfig config.Validator
}

// nolint: funlen
func (b Builder) Authenticators() map[string]Authenticator {
	all := make(map[string]Authenticator)

	for _, vendor := range b.Vendors {
		b.ValidateMappers(vendor.IssEntityMap, vendor.IssPeerMap)

		allowedAccessTypes := b.GetAllowedAccessTypes(vendor.AllowedAccessTypes)

		hid, err := topics.NewHashIDManager(vendor.HashIDMap)
		if err != nil {
			b.Logger.Fatal("cannot create hash-id manager", zap.Error(err))
		}

		var auth Authenticator

		switch {
		case vendor.UseValidator:
			client := validator.New(b.ValidatorConfig.URL, b.ValidatorConfig.Timeout)

			auth = &AutoAuthenticator{
				AllowedAccessTypes: allowedAccessTypes,
				Company:            vendor.Company,
				TopicManager: topics.NewTopicManager(
					vendor.Topics,
					hid,
					vendor.Company,
					vendor.IssEntityMap,
					vendor.IssPeerMap,
					b.Logger.Named("topic-manager"),
				),
				JwtConfig: vendor.Jwt,
				Validator: client,
				Parser:    jwt.NewParser(),
			}
		case vendor.IsInternal:
			if _, ok := vendor.Keys["system"]; !ok || len(vendor.Keys) != 1 {
				b.Logger.Fatal("admin authenticator supports only one key named system")
			}

			keys := b.GenerateKeys(vendor.Jwt.SigningMethod, vendor.Keys)

			auth = &AdminAuthenticator{
				Key:       keys["system"],
				Company:   vendor.Company,
				JwtConfig: vendor.Jwt,
				Parser:    jwt.NewParser(),
			}
		default:
			keys := b.GenerateKeys(vendor.Jwt.SigningMethod, vendor.Keys)

			auth = &ManualAuthenticator{
				Keys:               keys,
				AllowedAccessTypes: allowedAccessTypes,
				Company:            vendor.Company,
				TopicManager: topics.NewTopicManager(
					vendor.Topics,
					hid,
					vendor.Company,
					vendor.IssEntityMap,
					vendor.IssPeerMap,
					b.Logger.Named("topic-manager"),
				),
				JwtConfig: vendor.Jwt,
				Parser:    jwt.NewParser(jwt.WithValidMethods([]string{vendor.Jwt.SigningMethod})),
			}
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

// toUserAccessType will convert string access type to it's own type.
func toUserAccessType(access string) (acl.AccessType, error) {
	switch access {
	case "pub", "publish":
		return acl.Pub, nil
	case "sub", "subscribe":
		return acl.Sub, nil
	case "pubsub", "subpub":
		return acl.PubSub, nil
	}

	return "", ErrInvalidAccessType
}

func (b Builder) ValidateMappers(issEntityMap, issPeerMap map[string]string) {
	if _, ok := issEntityMap[topics.Default]; !ok {
		b.Logger.Fatal("default case for iss-entity map is required")
	}

	if _, ok := issPeerMap[topics.Default]; !ok {
		b.Logger.Fatal("default case for iss-peer map is required")
	}
}
