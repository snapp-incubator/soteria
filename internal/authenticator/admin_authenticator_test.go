package authenticator_test

import (
	"crypto/rsa"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AdminAuthenticatorTestSuite struct {
	suite.Suite

	AdminToken string
	key        *rsa.PrivateKey

	Authenticator authenticator.Authenticator
}

func TestAdminAuthenticator_suite(t *testing.T) {
	t.Parallel()

	st := new(AdminAuthenticatorTestSuite)

	pkey0, err := getPublicKey("admin")
	require.NoError(t, err)

	st.Authenticator = authenticator.AdminAuthenticator{
		Key:     pkey0,
		Company: "snapp-admin",
		Parser:  jwt.NewParser(),
		JwtConfig: config.JWT{
			IssName:       "iss",
			SubName:       "sub",
			SigningMethod: "rsa256",
		},
	}

	suite.Run(t, st)
}

func (suite *AdminAuthenticatorTestSuite) SetupSuite() {
	require := suite.Require()

	key, err := getPrivateKey("admin")
	require.NoError(err)

	suite.key = key

	adminToken, err := getSampleToken("admin", key)
	require.NoError(err)

	suite.AdminToken = adminToken
}

func (suite *AdminAuthenticatorTestSuite) TestAuth() {
	require := suite.Require()

	suite.Run("testing admin token auth", func() {
		require.NoError(suite.Authenticator.Auth(suite.AdminToken))
	})

	suite.Run("testing invalid token auth", func() {
		require.Error(suite.Authenticator.Auth(invalidToken))
	})

	suite.Run("testing invalid iss in auth token", func() {
		token, err := getSampleTokenWithClaims("admin", suite.key, "issuer", "sub")
		require.NoError(err)

		require.ErrorIs(suite.Authenticator.Auth(token), authenticator.ErrIssNotFound)
	})
}
