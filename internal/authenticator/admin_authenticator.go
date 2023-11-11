package authenticator

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/pkg/acl"
)

// AdminAuthenticator is responsible for Acl/Auth/Token of the internal system users,
// these users have admin access.
type AdminAuthenticator struct {
	Key       any
	Company   string
	JwtConfig config.Jwt
	Parser    *jwt.Parser
}

// Auth check user authentication by checking the user's token
// isSuperuser is a flag that authenticator set it true when credentials is related to a superuser.
func (a AdminAuthenticator) Auth(tokenString string) error {
	_, err := a.Parser.Parse(tokenString, func(
		token *jwt.Token,
	) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, ErrInvalidClaims
		}
		if claims[a.JwtConfig.IssName] == nil {
			return nil, ErrIssNotFound
		}

		return a.Key, nil
	})
	if err != nil {
		return fmt.Errorf("token is invalid: %w", err)
	}

	return nil
}

// ACL check a system user access to a topic.
// because we returns is-admin: true, this endpoint shouldn't
// be called.
func (a AdminAuthenticator) ACL(
	_ acl.AccessType,
	_ string,
	_ string,
) (bool, error) {
	return true, nil
}

func (a AdminAuthenticator) ValidateAccessType(_ acl.AccessType) bool {
	return true
}

func (a AdminAuthenticator) GetCompany() string {
	return a.Company
}
