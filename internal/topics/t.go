package topics

import "regexp"

type Topic string

type Type string

const (
	CabEvent         Type = "cab_event"
	DriverLocation   Type = "driver_location"
	SuperappEvent    Type = "superapp_event"
	BoxEvent         Type = "box_event"
	DaghighSys       Type = "daghigh_sys"
	GossiperLocation Type = "gossiper_location"
)

var (
	CabEventRegexp         = regexp.MustCompile(`(\w+)-event-[a-zA-Z0-9]+`)
	DriverLocationRegexp   = regexp.MustCompile(`snapp/driver/[a-zA-Z0-9+]+/location`)
	SuperappEventRegexp    = regexp.MustCompile(`snapp/(\w+)/[a-zA-Z0-9]+/superapp`)
	GossiperLocationRegexp = regexp.MustCompile(`snapp/passenger/[a-zA-Z0-9]+/driver-location`)
	DaghighSysRegexp       = regexp.MustCompile(`\$SYS/brokers/\+/clients/\+/(connected|disconnected)`)
)

func (t Topic) GetType() Type {
	topic := string(t)

	switch {
	case CabEventRegexp.MatchString(topic):
		return CabEvent
	case DriverLocationRegexp.MatchString(topic):
		return DriverLocation
	case SuperappEventRegexp.MatchString(topic):
		return SuperappEvent
	case topic == "bucks":
		return BoxEvent
	case GossiperLocationRegexp.MatchString(topic):
		return GossiperLocation
	case DaghighSysRegexp.MatchString(topic):
		return DaghighSys
	}

	return ""
}
