package user

import (
	"time"

	"github.com/google/uuid"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
)

// Type indicates user type which is herald, emq and staff.
type Type string

const (
	HeraldUser Type = "HeraldUser"
	EMQUser    Type = "EMQUser"
	Staff      Type = "Staff"
)

// Issuer indicate issuers.
type Issuer string

const (
	Driver     Issuer = "0"
	Passenger  Issuer = "1"
	ThirdParty Issuer = "100"
)

// User is Soteria's users db model.
type User struct {
	MetaData                db.MetaData   `json:"meta_data"`
	Username                string        `json:"username"`
	Password                string        `json:"password"`
	Type                    Type          `json:"type"`
	IPs                     []string      `json:"ips"`
	Secret                  string        `json:"secret"`
	TokenExpirationDuration time.Duration `json:"token_expiration_duration"`
	Rules                   []Rule        `json:"rules"`
}

// Rule tells about a access to a specific topic or endpoint.
type Rule struct {
	UUID       uuid.UUID      `json:"uuid"`
	Endpoint   string         `json:"endpoint"`
	Topic      topics.Type    `json:"topic"`
	AccessType acl.AccessType `json:"access_type"`
}

// GetMetadata is for getting metadata of a user model such as date created.
func (u User) GetMetadata() db.MetaData {
	return u.MetaData
}

// GetPrimaryKey is for knowing a model primary key.
func (u User) GetPrimaryKey() string {
	return u.Username
}

// CheckTopicAllowance checks whether the user is allowed to pub/sub/pubsub to a topic or not.
func (u User) CheckTopicAllowance(topic topics.Type, accessType acl.AccessType) bool {
	for _, rule := range u.Rules {
		if rule.Topic == topic && (rule.AccessType == acl.PubSub || rule.AccessType == accessType) {
			return true
		}
	}

	return false
}

// CheckEndpointAllowance checks whether the user is allowed to use a Herald endpoint or not.
func (u User) CheckEndpointAllowance(endpoint string) bool {
	for _, rule := range u.Rules {
		if rule.Endpoint == endpoint && rule.AccessType == acl.Pub {
			return true
		}
	}

	return false
}
