package authenticator

import (
	"github.com/snapp-incubator/soteria/pkg/acl"
)

type Authenticator interface {
	// Auth check user authentication by checking the user's token.
	// it returns error in case of any issue with the user token.
	Auth(tokenString string) error

	// ACL check a user access to a topic.
	ACL(
		accessType acl.AccessType,
		tokenString string,
		topic string,
	) (bool, error)

	// ValidateAccessType checks access type for specific topic
	ValidateAccessType(accessType acl.AccessType) bool

	// GetCompany Return the Company Field of The Inherited Objects
	GetCompany() string

	// IsSuperuser changes the Auth response in case of successful authentication
	// and shows user as superuser which disables the ACL.
	IsSuperuser() bool
}
