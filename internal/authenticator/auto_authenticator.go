package authenticator

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.snapp.ir/dispatching/soteria/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	validatorSDK "gitlab.snapp.ir/security_regulatory/validator/pkg/sdk"
)

// AutoAuthenticator is responsible for Acl/Auth/Token of users.
type AutoAuthenticator struct {
	Keys               map[string][]any
	AllowedAccessTypes []acl.AccessType
	TopicManager       *topics.Manager
	Company            string
	JwtConfig          config.Jwt
	Validator          validatorSDK.Client
}

// Auth check user authentication by checking the user's token
// isSuperuser is a flag that authenticator set it true when credentials is related to a superuser.
func (a AutoAuthenticator) Auth(tokenString string) error {
	if _, err := a.Validator.Validate(context.Background(), &http.Header{
		"X-Service-Name": []string{"Soteria"},
	}, tokenString); err != nil {
		return fmt.Errorf("token is invalid: %w", err)
	}

	return nil
}

// ACL check a user access to a topic.
// nolint: funlen, cyclop, dupl
func (a AutoAuthenticator) ACL(
	accessType acl.AccessType,
	tokenString string,
	topic string,
) (bool, error) {
	if !a.ValidateAccessType(accessType) {
		return false, ErrInvalidAccessType
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != a.JwtConfig.SigningMethod {
			return nil, ErrInvalidSigningMethod
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, ErrInvalidClaims
		}
		if claims[a.JwtConfig.IssName] == nil {
			return nil, ErrIssNotFound
		}
		if claims[a.JwtConfig.SubName] == nil {
			return nil, ErrSubNotFound
		}

		issuer := fmt.Sprintf("%v", claims[a.JwtConfig.IssName])
		key := a.Keys[issuer]
		if key == nil {
			return nil, KeyNotFoundError{Issuer: issuer}
		}

		return key, nil
	})
	if err != nil {
		return false, fmt.Errorf("token is invalid %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, ErrInvalidClaims
	}

	if claims[a.JwtConfig.IssName] == nil {
		return false, ErrIssNotFound
	}

	issuer := fmt.Sprintf("%v", claims[a.JwtConfig.IssName])

	if claims[a.JwtConfig.SubName] == nil {
		return false, ErrSubNotFound
	}

	sub, _ := claims[a.JwtConfig.SubName].(string)

	topicTemplate := a.TopicManager.ParseTopic(topic, issuer, sub)
	if topicTemplate == nil {
		return false, InvalidTopicError{Topic: topic}
	}

	if !topicTemplate.HasAccess(issuer, accessType) {
		return false, TopicNotAllowedError{
			issuer,
			sub,
			accessType,
			topic,
			topicTemplate.Type,
		}
	}

	return true, nil
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
