package authenticator

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/internal/topics"
	"github.com/snapp-incubator/soteria/pkg/acl"
	"github.com/snapp-incubator/soteria/pkg/validator"
)

// AutoAuthenticator is responsible for Acl/Auth/Token of users.
type AutoAuthenticator struct {
	AllowedAccessTypes []acl.AccessType
	TopicManager       *topics.Manager
	Company            string
	JwtConfig          config.Jwt
	Validator          validator.Client
	Parser             *jwt.Parser
}

// Auth check user authentication by checking the user's token
// isSuperuser is a flag that authenticator set it true when credentials is related to a superuser.
func (a AutoAuthenticator) Auth(tokenString string) error {
	if _, err := a.Validator.Validate(context.Background(), http.Header{
		validator.ServiceNameHeader: []string{"soteria"},
		"user-agent":                []string{},
		"X-APP-Version-Code":        []string{""},
		"X-APP-Version":             []string{""},
		"X-APP-Name":                []string{"soteria"},
		"locale":                    []string{"en-US"},
	}, fmt.Sprintf("bearer %s", tokenString)); err != nil {
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

	var claims jwt.MapClaims

	if _, _, err := a.Parser.ParseUnverified(tokenString, &claims); err != nil {
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

func (a AutoAuthenticator) IsSuperuser() bool {
	return false
}
