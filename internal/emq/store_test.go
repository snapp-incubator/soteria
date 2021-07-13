package emq_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/emq"
)

type StoreSuite struct {
	suite.Suite

	Client *redis.Client
	Store  emq.Store
}

func (suite *StoreSuite) SetupSuite() {
	mr, err := miniredis.Run()
	suite.Require().NoError(err)

	// nolint: exhaustivestruct
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	suite.Client = client
	suite.Store = emq.Store{Client: client}
}

func (suite *StoreSuite) TearDownTest() {
	suite.Require().NoError(suite.Client.FlushAll(context.Background()).Err())
}

func (suite *StoreSuite) TestSave() {
	require := suite.Require()

	ctx := context.Background()

	require.NoError(suite.Store.Save(ctx, emq.User{
		Username:    "secret",
		Password:    "password",
		IsSuperuser: true,
	}))

	ok, err := suite.Client.Exists(ctx, "mqtt_user:secret").Result()
	require.NoError(err)
	require.Equal(int64(1), ok)
}

func (suite *StoreSuite) TestLoad() {
	require := suite.Require()

	ctx := context.Background()

	user := emq.User{
		Username:    "secret",
		Password:    "password",
		IsSuperuser: true,
	}

	require.NoError(suite.Store.Save(ctx, user))

	u, err := suite.Store.Load(ctx, user.Username)
	require.NoError(err)
	require.Equal(u, user)
}

func (suite *StoreSuite) TestLoadAll() {
	require := suite.Require()

	total := 10

	ctx := context.Background()

	users := make([]emq.User, 0)

	for index := 0; index < total; index++ {
		user := emq.User{
			Username:    fmt.Sprintf("secret-%d", index),
			Password:    "password",
			IsSuperuser: true,
		}

		users = append(users, user)

		require.NoError(suite.Store.Save(ctx, user))
	}

	us, err := suite.Store.LoadAll(ctx)
	require.NoError(err)
	require.Equal(us, users)
}

func TestStoreSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(StoreSuite))
}
