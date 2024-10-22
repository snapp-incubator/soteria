package authenticator_test

import (
	"context"
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

	adminToken, err := getSampleToken("admin", key)
	require.NoError(err)

	suite.AdminToken = adminToken
}

func (suite *AdminAuthenticatorTestSuite) TestAuth() {
	require := suite.Require()

	suite.Run("testing admin token auth", func() {
		require.NoError(suite.Authenticator.Auth(context.Background(), suite.AdminToken))
	})

	suite.Run("testing invalid token auth", func() {
		require.Error(suite.Authenticator.Auth(context.Background(), invalidToken))
	})
}
