package authenticator

import (
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
)

type Authenticator interface {
	// Auth check user authentication by checking the user's token
	Auth(tokenString string) error

	// ACL check a user access to a topic.
	ACL(
		accessType acl.AccessType,
		tokenString string,
		topic string,
	) (bool, error)

	// GetCompany Return the Company Field of The Inherited Objects
	GetCompany() string
}
