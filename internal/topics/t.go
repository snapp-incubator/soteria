package topics

import (
	"crypto/md5"
	"errors"
	"fmt"
	"gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

const (
	CabEvent          string = "cab_event"
	DriverLocation           = "driver_location"
	PassengerLocation        = "passenger_location"
	SuperappEvent            = "superapp_event"
	BoxEvent                 = "box_event"
	SharedLocation           = "shared_location"
	Chat                     = "chat"
	GeneralCallEntry         = "general_call_entry"
	NodeCallEntry            = "node_call_entry"
	CallOutgoing             = "call_outgoing"
)

const (
	Driver    string = "driver"
	Passenger        = "passenger"
)

const (
	// EmqCabHashPrefix is the default prefix for hashing part of cab topic, default value is 'emqch'.
	EmqCabHashPrefix = "emqch"
	// EmqSuperAppHashPrefix is the default prefix for hashing part of super app topic, default value is 'superapp'.
	EmqSuperAppHashPrefix = "superapp"
)

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
			Regex:    regexp.MustCompile(topic.Regex),
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
	topicTemplate, ok := t.GetTopicTemplate(topic, "snapp")
	if !ok {
		return false
	}

	fields := make(map[string]string)
	audience := issuerToAudienceStr(issuer)
	fields["audience"] = audience
	fields["company"] = t.Company
	fields["peer"] = peerOfAudience(fields["audience"])

	hashId, err := t.getHashId(topicTemplate.Type, sub, issuer)
	if err != nil {
		return false
	}
	fields["hashId"] = hashId

	if topicTemplate.Type == NodeCallEntry {
		fields["node"] = strings.Split(topic, "/")[4]
	}

	parsedTopic := topicTemplate.Parse(fields)

	return parsedTopic == topic
}

func (t Manager) getHashId(topicType, sub string, issuer user.Issuer) (string, error) {
	switch topicType {
	case CabEvent, SuperappEvent:
		id, err := t.HashIDSManager.DecodeHashID(sub, issuerToAudience(issuer))
		if err != nil {
			return "", ErrDecodeHashID
		}

		prefix := EmqCabHashPrefix
		if topicType == SuperappEvent {
			prefix = EmqSuperAppHashPrefix
		}
		hid := md5.Sum([]byte(fmt.Sprintf("%s-%s", prefix, strconv.Itoa(id))))
		return fmt.Sprintf("%x", hid), nil
	default:
		return sub, nil
	}
}

func (t Manager) GetTopicTemplate(input, company string) (Template, bool) {
	topic := strings.TrimPrefix(input, company)

	for _, each := range t.TopicTemplates {
		if each.Regex.MatchString(topic) {
			return each, true
		}
	}

	return Template{}, false
}

// IsTopicValid returns true if it finds a topic type for the given topic.
func (t Manager) IsTopicValid(topic string) bool {
	return len(t.GetTopicType(topic, "snapp")) != 0
}

// GetTopicType finds topic type based on regexes.
func (t Manager) GetTopicType(input, company string) string {
	topic := strings.TrimPrefix(input, company)

	for _, each := range t.TopicTemplates {
		if each.Regex.MatchString(topic) {
			return each.Type
		}
	}

	return ""
}

// issuerToAudience returns corresponding audience in snappids form.
func issuerToAudience(issuer user.Issuer) snappids.Audience {
	switch issuer {
	case user.Passenger:
		return snappids.PassengerAudience
	case user.Driver:
		return snappids.DriverAudience
	default:
		return -1
	}
}

// issuerToAudienceStr returns corresponding audience in string form.
func issuerToAudienceStr(issuer user.Issuer) string {
	switch issuer {
	case user.Passenger:
		return Passenger
	case user.Driver:
		return Driver
	default:
		return ""
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
