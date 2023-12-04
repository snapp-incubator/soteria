package authenticator_test

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/internal/topics"
	"github.com/snapp-incubator/soteria/pkg/acl"
	"github.com/snapp-incubator/soteria/pkg/validator"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type AutoAuthenticatorTestSuite struct {
	suite.Suite

	Token     string
	PublicKey *rsa.PublicKey

	Server *httptest.Server

	Authenticator authenticator.Authenticator
}

func TestAutoAuthenticator_suite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(AutoAuthenticatorTestSuite))
}

func (suite *AutoAuthenticatorTestSuite) SetupSuite() {
	cfg := config.SnappVendor()

	require := suite.Require()

	pkey0, err := getPublicKey("0")
	require.NoError(err)

	suite.PublicKey = pkey0

	key0, err := getPrivateKey("0")
	require.NoError(err)

	token, err := getSampleToken("0", key0)
	require.NoError(err)

	suite.Token = token

	hid, err := topics.NewHashIDManager(cfg.HashIDMap)
	require.NoError(err)

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "bearer ")

		_, err := jwt.Parse(tokenString, func(
			token *jwt.Token,
		) (interface{}, error) {
			return pkey0, nil
		})
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)

			return
		}

		userData, err := json.Marshal(map[string]any{})
		require.NoError(err)
		res.Header().Add("X-User-Data", string(userData))

		res.WriteHeader(http.StatusOK)
	}))
	suite.Server = testServer

	suite.Authenticator = authenticator.AutoAuthenticator{
		Validator:          validator.New(testServer.URL, time.Second),
		AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub, acl.PubSub},
		Company:            "snapp",
		Parser:             jwt.NewParser(),
		TopicManager:       topics.NewTopicManager(cfg.Topics, hid, "snapp", cfg.IssEntityMap, cfg.IssPeerMap, zap.NewNop()),
		JWTConfig: config.JWT{
			IssName:       "iss",
			SubName:       "sub",
			SigningMethod: "rsa256",
		},
	}
}

func (suite *AutoAuthenticatorTestSuite) TestAuth() {
	require := suite.Require()

	suite.Run("testing valid token auth", func() {
		require.NoError(suite.Authenticator.Auth(suite.Token))
	})

	suite.Run("testing invalid token auth", func() {
		require.Error(suite.Authenticator.Auth(invalidToken))
	})
}

func (suite *AutoAuthenticatorTestSuite) TearDownSuite() {
	suite.Server.Close()
}

func TestAutoAuthenticator_ValidateTopicBySender(t *testing.T) {
	t.Parallel()

	cfg := config.SnappVendor()

	hid, err := topics.NewHashIDManager(cfg.HashIDMap)
	require.NoError(t, err)

	// nolint: exhaustruct
	authenticator := authenticator.AutoAuthenticator{
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
