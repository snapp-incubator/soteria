package authenticator

import (
	"gitlab.snapp.ir/dispatching/soteria/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
)

// AutoAuthenticator is responsible for Acl/Auth/Token of users.
type AutoAuthenticator struct {
	Keys               map[string][]any
	AllowedAccessTypes []acl.AccessType
	TopicManager       *topics.Manager
	Company            string
	JwtConfig          config.Jwt
}

// Auth check user authentication by checking the user's token
// isSuperuser is a flag that authenticator set it true when credentials is related to a superuser.
func (a AutoAuthenticator) Auth(tokenString string) error {
	_ = tokenString

	return nil
}

// ACL check a user access to a topic.
// nolint: funlen, cyclop
func (a AutoAuthenticator) ACL(
	accessType acl.AccessType,
	tokenString string,
	topic string,
) (bool, error) {
	_ = accessType
	_ = tokenString
	_ = topic

	return false, nil
}

func (a AutoAuthenticator) ValidateAccessType(accessType acl.AccessType) bool {
	for _, allowedAccessType := range a.AllowedAccessTypes {
		if allowedAccessType == accessType {
			return true
		}
	}

	return false
}

func (a AutoAuthenticator) GetCompany() string {
	return a.Company
}
