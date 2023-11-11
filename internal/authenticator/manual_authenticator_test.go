package authenticator_test

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/internal/topics"
	"github.com/snapp-incubator/soteria/pkg/acl"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

func TestManualAuthenticator_suite(t *testing.T) {
	t.Parallel()

	pkey0, err := getPublicKey("0")
	require.NoError(t, err)

	pkey1, err := getPublicKey("1")
	require.NoError(t, err)

	cfg := config.SnappVendor()

	hid, err := topics.NewHashIDManager(cfg.HashIDMap)
	require.NoError(t, err)

	st := new(AuthenticatorTestSuite)

	st.Authenticator = authenticator.ManualAuthenticator{
		Keys: map[string]any{
			topics.DriverIss:    pkey0,
			topics.PassengerIss: pkey1,
		},
		AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub, acl.PubSub},
		Company:            "snapp",
		Parser:             jwt.NewParser(),
		TopicManager:       topics.NewTopicManager(cfg.Topics, hid, "snapp", cfg.IssEntityMap, cfg.IssPeerMap, zap.NewNop()),
		JwtConfig: config.Jwt{
			IssName:       "iss",
			SubName:       "sub",
			SigningMethod: "rsa256",
		},
	}

	suite.Run(t, st)
}

func TestManualAuthenticator_ValidateTopicBySender(t *testing.T) {
	t.Parallel()

	cfg := config.SnappVendor()

	hid, err := topics.NewHashIDManager(cfg.HashIDMap)
	require.NoError(t, err)

	// nolint: exhaustruct
	authenticator := authenticator.ManualAuthenticator{
		AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub},
		Company:            "snapp",
		TopicManager:       topics.NewTopicManager(cfg.Topics, hid, "snapp", cfg.IssEntityMap, cfg.IssPeerMap, zap.NewNop()),
	}

	t.Run("testing valid driver cab event", func(t *testing.T) {
		t.Parallel()
		topicTemplate := authenticator.TopicManager.ParseTopic(validDriverCabEventTopic, topics.DriverIss, "DXKgaNQa7N5Y7bo")
		require.NotNil(t, topicTemplate)
	})
}
