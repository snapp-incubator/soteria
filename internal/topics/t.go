package topics

import (
	"regexp"
	"strings"
)

type Topic string

type Type string

const (
	CabEvent          Type = "cab_event"
	DriverLocation    Type = "driver_location"
	PassengerLocation Type = "passenger_location"
	SuperappEvent     Type = "superapp_event"
	BoxEvent          Type = "box_event"
	SharedLocation    Type = "shared_location"
	Chat              Type = "chat"
	CallEntry         Type = "call_entry"
	CallOutgoing      Type = "call_outgoing"
)

// Topic regular expressions which are used for detecting the topic name.
// topics are prefix with the company name will be trimed before matching
// so they regular expressions should not contain the company prefix.
var (
	CabEventRegexp          = regexp.MustCompile(`(\w+)-event-[a-zA-Z0-9]+`)
	DriverLocationRegexp    = regexp.MustCompile(`/driver/[a-zA-Z0-9+]+/location`)
	PassengerLocationRegexp = regexp.MustCompile(`/passenger/[a-zA-Z0-9+]+/location`)
	SuperappEventRegexp     = regexp.MustCompile(`/(driver|passenger)/[a-zA-Z0-9]+/(superapp)`)
	SharedLocationRegexp    = regexp.MustCompile(`/(driver|passenger)+/[a-zA-Z0-9]+/(driver-location|passenger-location)`)
	ChatRegexp              = regexp.MustCompile(`/(driver|passenger)+/[a-zA-Z0-9+]+/chat`)
	CallEntryRegexp         = regexp.MustCompile(`/(driver|passenger)+/[a-zA-Z0-9+]+/call/send`)
	CallOutgoingRegexp      = regexp.MustCompile(`/(driver|passenger)+/[a-zA-Z0-9+]+/call/receive`)
)

func (t Topic) GetType() Type {
	return t.GetTypeWithCompany("snapp")
}

func (t Topic) GetTypeWithCompany(company string) Type {
	topic := string(t)
	topic = strings.TrimPrefix(topic, company)

	switch {
	case CabEventRegexp.MatchString(topic):
		return CabEvent
	case DriverLocationRegexp.MatchString(topic):
		return DriverLocation
	case PassengerLocationRegexp.MatchString(topic):
		return PassengerLocation
	case SharedLocationRegexp.MatchString(topic):
		return SharedLocation
	case ChatRegexp.MatchString(topic):
		return Chat
	case CallEntryRegexp.MatchString(topic):
		return CallEntry
	case CallOutgoingRegexp.MatchString(topic):
		return CallOutgoing
	case SuperappEventRegexp.MatchString(topic):
		return SuperappEvent
	case topic == "bucks":
		return BoxEvent
	}

	return ""
}
