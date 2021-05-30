package topics

import "regexp"

type Topic string

type Type string

const (
	CabEvent          Type = "cab_event"
	DriverLocation    Type = "driver_location"
	PassengerLocation Type = "passenger_location"
	SuperappEvent     Type = "superapp_event"
	BoxEvent          Type = "box_event"
	DaghighSys        Type = "daghigh_sys"
	SharedLocation    Type = "shared_location"
	Chat              Type = "chat"
)

var (
	CabEventRegexp          = regexp.MustCompile(`(\w+)-event-[a-zA-Z0-9]+`)
	DriverLocationRegexp    = regexp.MustCompile(`snapp/driver/[a-zA-Z0-9+]+/location`)
	PassengerLocationRegexp = regexp.MustCompile(`snapp/passenger/[a-zA-Z0-9+]+/location`)
	SuperappEventRegexp     = regexp.MustCompile(`snapp/(driver|passenger)/[a-zA-Z0-9]+/(superapp)`)
	SharedLocationRegexp    = regexp.MustCompile(`snapp+/(driver|passenger)+/[a-zA-Z0-9]+/(driver-location|passenger-location)`)
	DaghighSysRegexp        = regexp.MustCompile(`\$SYS/brokers/\+/clients/\+/(connected|disconnected)`)
	ChatRegexp              = regexp.MustCompile(`snapp+/(driver|passenger)+/[a-zA-Z0-9]+/(driver-chat|passenger-chat)`)
)

func (t Topic) GetType() Type {
	topic := string(t)

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
	case SuperappEventRegexp.MatchString(topic):
		return SuperappEvent
	case topic == "bucks":
		return BoxEvent
	case DaghighSysRegexp.MatchString(topic):
		return DaghighSys
	}

	return ""
}
