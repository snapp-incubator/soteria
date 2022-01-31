package topics

import (
	"crypto/md5" // nolint:gosec
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
)

const (
	CabEvent          string = "cab_event"
	DriverLocation    string = "driver_location"
	PassengerLocation string = "passenger_location"
	SuperappEvent     string = "superapp_event"
	BoxEvent          string = "box_event"
	SharedLocation    string = "shared_location"
	Chat              string = "chat"
	GeneralCallEntry  string = "general_call_entry"
	NodeCallEntry     string = "node_call_entry"
	CallOutgoing      string = "call_outgoing"
)

const (
	Driver    string = "driver"
	Passenger string = "passenger"
)

// EmqCabHashPrefix is the default prefix for hashing part of cab topic, default value is 'emqch'.
const EmqCabHashPrefix = "emqch"

var ErrDecodeHashID = errors.New("could not decode hash id")

// Topic regular expressions which are used for detecting the topic name.
// topics are prefix with the company name will be trimed before matching
// so they regular expressions should not contain the company prefix.

type Manager struct {
	HashIDSManager *snappids.HashIDSManager
	Company        string
	TopicTemplates []Template
}

// NewTopicManager returns a topic manager to validate topics.
func NewTopicManager(topicList []Topic, hashIDManager *snappids.HashIDSManager, company string) Manager {
	templates := make([]Template, 0)

	for _, topic := range topicList {
		each := Template{
			Type:     topic.Type,
			Template: template.Must(template.New(topic.Type).Parse(topic.Template)),
			HashType: topic.HashType,
			Accesses: topic.Accesses,
		}
		templates = append(templates, each)
	}

	return Manager{
		HashIDSManager: hashIDManager,
		Company:        company,
		TopicTemplates: templates,
	}
}

func (t Manager) ValidateTopicBySender(topic string, issuer user.Issuer, sub string) bool {
	topicTemplate, ok := t.GetTopicTemplate(topic)
	if !ok {
		return false
	}

	fields := make(map[string]string)
	audience := issuerToAudienceStr(issuer)
	fields["audience"] = audience
	fields["company"] = t.Company
	fields["peer"] = peerOfAudience(fields["audience"])

	hashID, err := t.getHashID(topicTemplate.Type, sub, issuer)
	if err != nil {
		return false
	}

	fields["hashId"] = hashID

	if topicTemplate.Type == NodeCallEntry {
		fields["node"] = strings.Split(topic, "/")[4]
	}

	parsedTopic := topicTemplate.Parse(fields)

	return parsedTopic == topic
}

func (t Manager) getHashID(hashType HashType, sub string, audience snappids.Audience) (string, error) {
	if hashType == MD5 {
		id, err := t.HashIDSManager.DecodeHashID(sub, audience)
		if err != nil {
			return "", ErrDecodeHashID
		}

		hid := md5.Sum([]byte(fmt.Sprintf("%s-%s", EmqCabHashPrefix, strconv.Itoa(id)))) // nolint:gosec

		return fmt.Sprintf("%x", hid), nil
	}

	return sub, nil
}

// IssuerToAudience returns corresponding audience in snappids form.
func IssuerToAudience(issuer user.Issuer) (snappids.Audience, string) {
	switch issuer {
	case user.Passenger:
		return snappids.PassengerAudience, Passenger
	case user.Driver:
		return snappids.DriverAudience, Driver
	default:
		return -1, ""
	}
}

func peerOfAudience(audience string) string {
	switch audience {
	case Driver:
		return Passenger
	case Passenger:
		return Driver
	default:
		return ""
	}
}
