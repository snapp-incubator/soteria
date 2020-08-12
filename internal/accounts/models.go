package accounts

import (
	"crypto/rsa"
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"regexp"
	"time"
)

var ModelHandler db.ModelHandler

const (
	// User Types
	HeraldUser = "HeraldUser"
	EMQUser    = "EMQUser"
	Staff      = "Staff"
)

// User is Soteria's users db model
type User struct {
	MetaData                db.MetaData    `json:"meta_data"`
	Username                string         `json:"username"`
	Password                []byte         `json:"password"`
	Type                    string         `json:"type"`
	IPs                     []string         `json:"ips"`
	Secret                  string         `json:"secret"`
	PublicKey               *rsa.PublicKey `json:"public_key"`
	TokenExpirationDuration time.Duration  `json:"token_expiration_duration"`
	Rules                   []Rule         `json:"rules"`
}

// Rule tells about a access to a specific topic or endpoint
type Rule struct {
	UID          int    `json:"uid"`
	Endpoint     string `json:"endpoint"`
	TopicPattern string `json:"topic_pattern"`
	AccessType   string `json:"access_type"`
}

// GetMetadata is for getting metadata of a user model such as date created
func (u User) GetMetadata() db.MetaData {
	return u.MetaData
}

// GetPrimaryKey is for knowing a model primary key
func (u User) GetPrimaryKey() string {
	return u.Username
}

// CheckTopicAllowance checks whether the user is allowed to pub/sub/pubsub to a topic or not
func (u User) CheckTopicAllowance(topicMan TopicMan, issuer, sub, topic, accessType string) bool {
	for _, rule := range u.Rules {
		matched, _ := regexp.Match(rule.TopicPattern, []byte(topic))
		if rule.AccessType == accessType && matched && topicMan(issuer, sub) == topic {
			return true
		}
	}
	return false
}

// CheckEndpointAllowance checks whether the user is allowed to use a Herald endpoint or not
func (u User) CheckEndpointAllowance(endpoint string) bool {
	for _, rule := range u.Rules {
		if rule.Endpoint == endpoint && rule.AccessType == Pub {
			return true
		}
	}
	return false
}
