// Package topics regular expressions which are used for detecting the topic name.
// topics are prefix with the company name will be trimed before matching
// so they regular expressions should not contain the company prefix.
package topics

import (
	"crypto/md5" // nolint: gosec
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
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
	DriverIss string = "0"

	Passenger    string = "passenger"
	PassengerIss string = "1"
)

// EmqCabHashPrefix is the default prefix for hashing part of cab topic, default value is 'emqch'.
const EmqCabHashPrefix = "emqch"

var ErrDecodeHashID = errors.New("could not decode hash id")

type Manager struct {
	HashIDSManager *snappids.HashIDSManager
	Company        string
	TopicTemplates []Template
	IssEntityMap   map[string]string
	IssPeerMap     map[string]string
	Functions      template.FuncMap
}

// NewTopicManager returns a topic manager to validate topics.
func NewTopicManager(topicList []Topic, hashIDManager *snappids.HashIDSManager, company string, issEntityMap, issPeerMap map[string]string) *Manager {
	manager := &Manager{
		HashIDSManager: hashIDManager,
		Company:        company,
		IssEntityMap:   issEntityMap,
		IssPeerMap:     issPeerMap,
	}

	manager.Functions = template.FuncMap{
		"IssToEntity":  manager.IssEntityMapper,
		"HashID":       manager.getHashID,
		"IssToSnappID": manager.IssToSnappID,
		"IssToPeer":    manager.IssPeerMapper,
	}

	templates := make([]Template, 0)

	for _, topic := range topicList {
		each := Template{
			Type:     topic.Type,
			Template: template.Must(template.New(topic.Type).Funcs(manager.Functions).Parse(topic.Template)),
			HashType: topic.HashType,
			Accesses: topic.Accesses,
		}
		templates = append(templates, each)
	}

	manager.TopicTemplates = templates

	return manager
}

// ValidateTopic checks if a topic is valid based on the given parameters.
func (t Manager) ValidateTopic(topic, iss, sub string) *Template {
	fields := make(map[string]any)
	fields["iss"] = iss
	fields["company"] = t.Company
	fields["sub"] = sub

	for _, topicTemplate := range t.TopicTemplates {
		fields["hashType"] = topicTemplate.HashType

		regex := new(strings.Builder)

		err := topicTemplate.Template.Execute(regex, fields)
		if err != nil {
			return nil
		}

		if regexp.MustCompile(regex.String()).MatchString(topic) {
			return &topicTemplate
		}
	}

	return nil
}

// getHashID calculate hashID based on hashType.
// most of the topics have hashID type for their hashIDs but some old topics have different hashTypes.
// if hashType is equal to hashID, sub is returned without any changes.
func (t Manager) getHashID(hashType HashType, sub string, audience snappids.Audience) string {
	if hashType == MD5 {
		id, err := t.HashIDSManager.DecodeHashID(sub, audience)
		if err != nil {
			return ""
		}

		hid := md5.Sum([]byte(fmt.Sprintf("%s-%s", EmqCabHashPrefix, strconv.Itoa(id)))) // nolint:gosec

		return fmt.Sprintf("%x", hid)
	}

	return sub
}

func (t *Manager) IssEntityMapper(iss string) string {
	result, ok := t.IssEntityMap[iss]
	if ok {
		return result
	}

	return iss
}

func (t *Manager) IssPeerMapper(iss string) string {
	result, ok := t.IssPeerMap[iss]
	if ok {
		return result
	}

	return iss
}

// IssuerToAudience returns corresponding audience in snappids form.
func IssuerToAudience(issuer string) (snappids.Audience, string) {
	switch issuer {
	case user.Passenger:
		return snappids.PassengerAudience, Passenger
	case user.Driver:
		return snappids.DriverAudience, Driver
	case user.None:
		fallthrough
	default:
		return -1, ""
	}
}

// IssToSnappID returns corresponding audience in snappids form.
func (t *Manager) IssToSnappID(iss string) snappids.Audience {
	switch iss {
	case user.Passenger:
		return snappids.PassengerAudience
	case user.Driver:
		return snappids.DriverAudience
	case user.None:
		fallthrough
	default:
		return -1
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
