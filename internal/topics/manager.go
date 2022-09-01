// Package topics regular expressions which are used for detecting the topic name.
// topics are prefix with the company name will be trimed before matching
// so they regular expressions should not contain the company prefix.
package topics

import (
	"crypto/md5" //nolint: gosec
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/speps/go-hashids/v2"
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

	NoneIss string = "-1"
)

const (
	// EmqCabHashPrefix is the default prefix for hashing part of cab topic, default value is 'emqch'.
	EmqCabHashPrefix = "emqch"

	Default = "default"
)

type Manager struct {
	HashIDSManager map[string]*hashids.HashID
	Company        string
	TopicTemplates []Template
	IssEntityMap   map[string]string
	IssPeerMap     map[string]string
	Functions      template.FuncMap
}

// NewTopicManager returns a topic manager to validate topics.
func NewTopicManager(
	topicList []Topic,
	hashIDManager map[string]*hashids.HashID,
	company string,
	issEntityMap, issPeerMap map[string]string,
) *Manager {
	manager := &Manager{ //nolint: exhaustruct
		HashIDSManager: hashIDManager,
		Company:        company,
		IssEntityMap:   issEntityMap,
		IssPeerMap:     issPeerMap,
	}

	manager.Functions = template.FuncMap{
		"IssToEntity": manager.IssEntityMapper,
		"HashID":      manager.getHashID,
		"IssToPeer":   manager.IssPeerMapper,
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

// ParseTopic checks if a topic is valid based on the given parameters.
func (t *Manager) ParseTopic(topic, iss, sub string) *Template {
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
func (t *Manager) getHashID(hashType HashType, sub string, iss string) string {
	if hashType == MD5 {
		id, err := t.HashIDSManager[iss].DecodeWithError(sub)
		if err != nil {
			return ""
		}

		hid := md5.Sum([]byte(fmt.Sprintf("%s-%s", EmqCabHashPrefix, strconv.Itoa(id[0])))) //nolint:gosec

		return fmt.Sprintf("%x", hid)
	}

	return sub
}

func (t *Manager) IssEntityMapper(iss string) string {
	result, ok := t.IssEntityMap[iss]
	if ok {
		return result
	}

	return t.IssEntityMap[Default]
}

func (t *Manager) IssPeerMapper(iss string) string {
	result, ok := t.IssPeerMap[iss]
	if ok {
		return result
	}

	return t.IssPeerMap[Default]
}
