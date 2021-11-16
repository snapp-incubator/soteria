package acl

// Access Types for EMQ contains subscribe, publish and publish-subscribe.
type AccessType string

const (
	Sub    AccessType = "1"
	Pub    AccessType = "2"
	PubSub AccessType = "3"

	ClientCredentials = "client_credentials"
)

func (a AccessType) String() string {
	switch a {
	case Sub:
		return "subscribe"
	case Pub:
		return "publish"
	case PubSub:
		return "publish-subscribe"
	}

	return ""
}
