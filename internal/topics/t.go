package topics

import "regexp"

type Topic string

type Type string

const (
	CabEvent       Type = "cab_event"
	DriverLocation Type = "driver_location"
	SuperappEvent  Type = "superapp_event"
	BoxEvent       Type = "box_event"
)

func (t Topic) GetType() Type {
	matched, _ := regexp.Match(`(\w+)-event-[a-zA-Z0-9]+`, []byte(t))
	if matched {
		return CabEvent
	}
	matched, _ = regexp.Match(`snapp/driver/[a-zA-Z0-9]+/location`, []byte(t))
	if matched {
		return DriverLocation
	}
	matched, _ = regexp.Match(`snapp/(\w+)/[a-zA-Z0-9]+/superapp`, []byte(t))
	if matched {
		return SuperappEvent
	}
	if t == "bucks" {
		return BoxEvent
	}
	return ""
}
