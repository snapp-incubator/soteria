package user

import (
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
)

// Issuer indicate issuers.
type Issuer string

const (
	Driver    Issuer = "0"
	Passenger Issuer = "1"
)

// User is Soteria's users db model.
type User struct {
	Username string `json:"username"`
	Rules    []Rule `json:"rules"`
}

// Rule tells about a access to a specific topic or endpoint.
type Rule struct {
	Topic      topics.Type    `json:"topic"`
	AccessType acl.AccessType `json:"access_type"`
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
