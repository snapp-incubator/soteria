package authenticator

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/internal/metric"
	"github.com/snapp-incubator/soteria/internal/topics"
	"github.com/snapp-incubator/soteria/pkg/acl"
	"github.com/snapp-incubator/soteria/pkg/validator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// AutoAuthenticator is responsible for Acl/Auth/Token of users.
type AutoAuthenticator struct {
	AllowedAccessTypes []acl.AccessType
	TopicManager       *topics.Manager
	Company            string
	JWTConfig          config.JWT
	Validator          validator.Client
	Parser             *jwt.Parser
	Tracer             trace.Tracer
	metrics            *metric.AutoAuthenticatorMetrics
}

// Auth check user authentication by checking the user's token
// isSuperuser is a flag that authenticator set it true when credentials is related to a superuser.
func (a AutoAuthenticator) Auth(ctx context.Context, tokenString string) error {
	ctx, span := a.Tracer.Start(ctx, "auto-authenticator.auth")
	span.End()

	headers := http.Header{
		validator.ServiceNameHeader: []string{"soteria"},
		"user-agent":                []string{},
		"X-APP-Version-Code":        []string{""},
		"X-APP-Version":             []string{""},
		"X-APP-Name":                []string{"soteria"},
		"X-Original-URI":            []string{"/v2/auth"},
		"locale":                    []string{"en-US"},
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))

	start := time.Now()

	if _, err := a.Validator.Validate(ctx, headers, "bearer "+tokenString); err != nil {
		a.metrics.Latency(time.Since(start).Seconds(), a.Company, err)

		return fmt.Errorf("token is invalid: %w", err)
	}

	a.metrics.Latency(time.Since(start).Seconds(), a.Company, nil)

	return nil
}

// ACL check a user access to a topic.
// nolint: cyclop, dupl
func (a AutoAuthenticator) ACL(
	_ context.Context,
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

	if claims[a.JWTConfig.IssName] == nil {
		return false, ErrIssNotFound
	}

	issuer := fmt.Sprintf("%v", claims[a.JWTConfig.IssName])

	if claims[a.JWTConfig.SubName] == nil {
		return false, ErrSubNotFound
	}

	sub := fmt.Sprintf("%v", claims[a.JWTConfig.SubName])

	topicTemplate := a.TopicManager.ParseTopic(topic, issuer, sub, map[string]any(claims))
	if topicTemplate == nil {
		return false, InvalidTopicError{Topic: topic}
	}

	if !topicTemplate.HasAccess(issuer, accessType) {
		return false, TopicNotAllowedError{
			Issuer:     issuer,
			Sub:        sub,
			AccessType: accessType,
			Topic:      topic,
			TopicType:  topicTemplate.Type,
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
