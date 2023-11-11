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

	st := new(AuthenticatorTestSuite)

	pkey0, err := st.getPublicKey(topics.DriverIss)
	require.NoError(t, err)

	pkey1, err := st.getPublicKey(topics.PassengerIss)
	require.NoError(t, err)

	cfg := config.SnappVendor()

	hid, err := topics.NewHashIDManager(cfg.HashIDMap)
	require.NoError(t, err)

	// nolint: exhaustruct
	st.Authenticator = authenticator.ManualAuthenticator{
		Keys: map[string]any{
			topics.DriverIss:    pkey0,
			topics.PassengerIss: pkey1,
		},
		AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub, acl.PubSub},
		Company:            "snapp",
		Parser:             jwt.NewParser(),
		TopicManager:       topics.NewTopicManager(cfg.Topics, hid, "snapp", cfg.IssEntityMap, cfg.IssPeerMap, zap.NewNop()),
	}

	suite.Run(t, st)
}
