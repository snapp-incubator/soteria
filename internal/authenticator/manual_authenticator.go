package authenticator

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/internal/topics"
	"github.com/snapp-incubator/soteria/pkg/acl"
)

// ManualAuthenticator is responsible for Acl/Auth/Token of userIDs without calling
// any http client, etc.
type ManualAuthenticator struct {
	Keys               map[string]any
	AllowedAccessTypes []acl.AccessType
	TopicManager       *topics.Manager
	Company            string
	JWTConfig          config.JWT
	Parser             *jwt.Parser
}

// Auth check user authentication by checking the user's token.
func (a ManualAuthenticator) Auth(tokenString string) error {
	_, err := a.Parser.Parse(tokenString, func(
		token *jwt.Token,
	) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, ErrInvalidClaims
		}

		if claims[a.JWTConfig.IssName] == nil {
			return nil, ErrIssNotFound
		}

		issuer := fmt.Sprintf("%v", claims[a.JWTConfig.IssName])

		key := a.Keys[issuer]
		if key == nil {
			return nil, KeyNotFoundError{Issuer: issuer}
		}

		return key, nil
	})
	if err != nil {
		return fmt.Errorf("token is invalid: %w", err)
	}

	return nil
}

// ACL check a user access to a topic.
// nolint: funlen, cyclop, dupl
func (a ManualAuthenticator) ACL(
	accessType acl.AccessType,
	tokenString string,
	topic string,
) (bool, error) {
	if !a.ValidateAccessType(accessType) {
		return false, ErrInvalidAccessType
	}

	token, err := a.Parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, ErrInvalidClaims
		}

		if claims[a.JWTConfig.IssName] == nil {
			return nil, ErrIssNotFound
		}

		if claims[a.JWTConfig.SubName] == nil {
			return nil, ErrSubNotFound
		}

		issuer := fmt.Sprintf("%v", claims[a.JWTConfig.IssName])

		key := a.Keys[issuer]
		if key == nil {
			return nil, KeyNotFoundError{Issuer: issuer}
		}

		return key, nil
	})
	if err != nil {
		return false, fmt.Errorf("token is invalid: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, ErrInvalidClaims
	}

	if claims[a.JWTConfig.IssName] == nil {
		return false, ErrIssNotFound
	}

	issuer := fmt.Sprintf("%v", claims[a.JWTConfig.IssName])

	if claims[a.JWTConfig.SubName] == nil {
		return false, ErrSubNotFound
	}

	sub, _ := claims[a.JWTConfig.SubName].(string)

	topicTemplate := a.TopicManager.ParseTopic(topic, issuer, sub, map[string]any(claims))
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

func (a ManualAuthenticator) ValidateAccessType(accessType acl.AccessType) bool {
	for _, allowedAccessType := range a.AllowedAccessTypes {
		if allowedAccessType == accessType {
			return true
		}
	}

	return false
}

func (a ManualAuthenticator) GetCompany() string {
	return a.Company
}

func (a ManualAuthenticator) IsSuperuser() bool {
	return false
}
