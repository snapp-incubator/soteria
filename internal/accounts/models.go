package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"net"
)

var ModelHandler db.ModelHandler

const (
	HERALD = "HERALD"
	EMQ    = "EMQ"
	STAFF  = "STAFF"
)

type User struct {
	MetaData              db.MetaData `json:"meta_data"`
	Username              string      `json:"username"`
	Password              []byte      `json:"password"`
	Type                  string      `json:"type"`
	IPs                   []net.IP    `json:"ips"`
	AuthorizedEndpoints   string      `json:"authorized_endpoints"`
	UnauthorizedEndpoints string      `json:"unauthorized_endpoints"`
	Secret                string      `json:"secret"`
	AccessType            string      `json:"access_type"`
	AuthorizedTopics      string      `json:"authorized_topics"`
}

func (u User) GetMetadata() db.MetaData {
	return u.MetaData
}

func (u User) GetPrimaryKey() string {
	return u.Username
}
