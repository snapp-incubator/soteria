package cachedredis_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/suite"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db/cachedredis"
	rd "gitlab.snapp.ir/dispatching/soteria/v3/internal/db/redis"
)

type CacheModelHandlerSuite struct {
	suite.Suite

	Model db.ModelHandler
	DB    *redis.Client
	Cache *cache.Cache
}

func (suite *CacheModelHandlerSuite) SetupSuite() {
	mr, err := miniredis.Run()
	suite.NoError(err)

	// nolint: exhaustivestruct
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	cache := cache.New(time.Minute*60, time.Minute*60)

	suite.DB = client
	suite.Cache = cache

	suite.Model = cachedredis.NewCachedRedisModelHandler(&rd.ModelHandler{Client: client}, cache)
}

func (suite *CacheModelHandlerSuite) TestGet() {
	require := suite.Require()

	expected := MockModel{Name: "test", Value: ""}
	v, err := json.Marshal(expected)
	require.NoError(err)

	require.NoError(suite.DB.Set(context.Background(), "mock-test", v, 0).Err())

	suite.Run("testing first successful get", func() {
		var actual MockModel
		require.NoError(suite.Model.Get(context.Background(), "mock", "test", &actual))
		require.Equal(expected, actual)

		suite.Cache.Get(rd.GenerateKey("mock", "test"))
	})
	suite.Run("testing second successful get", func() {
		var actual MockModel
		err := suite.Model.Get(context.Background(), "mock", "test", &actual)

		require.Equal(expected, actual)
		require.NoError(err)
	})
	suite.Run("testing failed get", func() {
		var actual MockModel
		err := suite.Model.Get(context.Background(), "mock", "t", &actual)

		require.EqualError(err, "redis: nil")
	})

	suite.NoError(suite.DB.Del(context.Background(), "mock-test").Err())
}

func (suite *CacheModelHandlerSuite) TestSave() {
	require := suite.Require()

	expected := MockModel{Name: "save-test", Value: ""}

	require.NoError(suite.Model.Save(context.Background(), expected))

	v := suite.DB.Get(context.Background(), "mock-save-test").Val()

	var actual MockModel

	require.NoError(json.Unmarshal([]byte(v), &actual))
	require.Equal(expected, actual)
}

func (suite *CacheModelHandlerSuite) TestDelete() {
	require := suite.Require()

	expected := MockModel{Name: "test", Value: ""}
	v, err := json.Marshal(expected)
	require.NoError(err)

	require.NoError(suite.DB.Set(context.Background(), "mock-test", v, 0).Err())

	suite.Run("testing successful delete", func() {
		require.NoError(suite.Model.Delete(context.Background(), "mock", "test"))

		err := suite.DB.Get(context.Background(), "mock-test").Err()
		require.EqualError(err, "redis: nil")

		item, ok := suite.Cache.Get("mock-test")
		require.False(ok)
		require.Nil(item)
	})
	suite.Run("testing failed delete", func() {
		var actual MockModel

		err := suite.Model.Get(context.Background(), "mock", "t", &actual)
		require.EqualError(err, "redis: nil")
	})
}

func (suite *CacheModelHandlerSuite) TestUpdate() {
	require := suite.Require()

	m := MockModel{Name: "test", Value: "test-1"}
	v, err := json.Marshal(m)
	require.NoError(err)

	require.NoError(suite.DB.Set(context.Background(), "mock-test", v, 0).Err())

	newModel := MockModel{Name: "test", Value: "test-2"}
	require.NoError(suite.Model.Update(context.Background(), newModel))

	var updatedModel MockModel

	require.NoError(suite.Model.Get(context.Background(), "mock", "test", &updatedModel))
	require.Equal(newModel, updatedModel)
}

type MockModel struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (m MockModel) GetMetadata() db.MetaData {
	return db.MetaData{
		ModelName:    "mock",
		DateCreated:  time.Time{},
		DateModified: time.Time{},
	}
}

func (m MockModel) GetPrimaryKey() string {
	return m.Name
}

func TestCacheModelHandlerSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(CacheModelHandlerSuite))
}
