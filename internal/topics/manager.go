// Package topics regular expressions which are used for detecting the topic name.
// topics are prefix with the company name will be trimemd before matching
// so they regular expressions should not contain the company prefix.
package topics

import (
	"crypto/md5" //nolint: gosec
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	jwtstrconv "github.com/snapp-incubator/soteria/pkg/strconv"
	"github.com/speps/go-hashids/v2"
	regexp "github.com/wasilibs/go-re2"
	"go.uber.org/zap"
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
	Logger         *zap.Logger
}

// NewTopicManager returns a topic manager to validate topics.
func NewTopicManager(
	topicList []Topic,
	hashIDManager map[string]*hashids.HashID,
	company string,
	issEntityMap, issPeerMap map[string]string,
	logger *zap.Logger,
) *Manager {
	manager := &Manager{ //nolint: exhaustruct
		HashIDSManager: hashIDManager,
		Company:        company,
		IssEntityMap:   issEntityMap,
		IssPeerMap:     issPeerMap,
		Logger: logger.With(
			zap.String("company", company),
		),
	}

	manager.Functions = template.FuncMap{
		"IssToEntity":  manager.IssEntityMapper,
		"DecodeHashID": manager.DecodeHashID,
		"EncodeHashID": manager.EncodeHashID,
		"EncodeMD5":    manager.EncodeMD5,
		"IssToPeer":    manager.IssPeerMapper,
	}

	templates := make([]Template, 0, len(topicList))

	for _, topic := range topicList {
		each := Template{
			Type:     topic.Type,
			Template: template.Must(template.New(topic.Type).Funcs(manager.Functions).Parse(topic.Template)),
			Accesses: topic.Accesses,
		}
		templates = append(templates, each)
	}

	manager.TopicTemplates = templates

	return manager
}

// ParseTopic checks if a topic is valid based on the given parameters.
func (t *Manager) ParseTopic(topic, iss, sub string, claims map[string]any) *Template {
	fields := make(map[string]string)

	for k, v := range claims {
		fields[k] = jwtstrconv.ToString(v)
	}

	fields["iss"] = iss
	fields["company"] = t.Company
	fields["sub"] = sub

	for _, topicTemplate := range t.TopicTemplates {
		regex := new(strings.Builder)

		if err := topicTemplate.Template.Execute(regex, fields); err != nil {
			t.Logger.Error("template execution failed", zap.Error(err), zap.String("template", topicTemplate.Type))

			return nil
		}

		t.Logger.Debug("topic template generated",
			zap.String("topic", regex.String()),
			zap.String("iss", iss),
			zap.String("sub", sub),
		)

		if regexp.MustCompile(regex.String()).MatchString(topic) {
			return &topicTemplate
		}
	}

	return nil
}

func (t *Manager) EncodeMD5(iss string) string {
	hid := md5.Sum(fmt.Appendf(nil, "%s-%s", EmqCabHashPrefix, iss)) //nolint:gosec

	return hex.EncodeToString(hid[:])
}

func (t *Manager) DecodeHashID(sub, iss string) string {
	id, err := t.HashIDSManager[iss].DecodeWithError(sub)
	if err != nil {
		t.Logger.Error("decoding sub failed", zap.Error(err), zap.String("sub", sub))

		return ""
	}

	return strconv.Itoa(id[0])
}

func (t *Manager) EncodeHashID(sub, iss string) string {
	subInt, err := strconv.Atoi(sub)
	if err != nil {
		t.Logger.Error("encoding sub failed", zap.Error(err), zap.String("sub", sub))

		return ""
	}

	id, err := t.HashIDSManager[iss].Encode([]int{subInt})
	if err != nil {
		t.Logger.Error("encoding sub failed", zap.Error(err), zap.String("sub", sub))

		return ""
	}

	return id
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

func NewHashIDManager(hidmap map[string]HashData) (map[string]*hashids.HashID, error) {
	hid := make(map[string]*hashids.HashID)

	for iss, data := range hidmap {
		var err error

		hd := hashids.NewData()
		hd.Salt = data.Salt
		hd.MinLength = data.Length

		if data.Alphabet != "" {
			hd.Alphabet = data.Alphabet
		}

		hid[iss], err = hashids.NewWithData(hd)
		if err != nil {
			return nil, fmt.Errorf("cannot create hashid enc/dec %w", err)
		}
	}

	return hid, nil
}
