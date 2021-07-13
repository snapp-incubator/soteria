package redis_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	rd "gitlab.snapp.ir/dispatching/soteria/v3/internal/db/redis"
)

type RedisModelHandlerSuite struct {
	suite.Suite

	Model db.ModelHandler
	DB    *redis.Client
}

func (suite *RedisModelHandlerSuite) SetupSuite() {
	mr, err := miniredis.Run()
	suite.Require().NoError(err)

	// nolint: exhaustivestruct
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	suite.DB = client
	suite.Model = rd.ModelHandler{Client: client}
}

func (suite *RedisModelHandlerSuite) TestGet() {
	require := suite.Require()

	expected := MockModel{Name: "test", Value: ""}
	v, err := json.Marshal(expected)
	require.NoError(err)

	require.NoError(suite.DB.Set(context.Background(), "mock-test", v, 0).Err())

	suite.Run("testing successful get", func() {
		var actual MockModel
		err := suite.Model.Get(context.Background(), "mock", "test", &actual)
		require.Equal(expected, actual)
		require.NoError(err)
	})

	suite.Run("testing invalid get usage", func() {
		var actual MockModel
		err := suite.Model.Get(context.Background(), "mock", "test", actual)
		require.Error(err)
	})

	suite.Run("testing failed get", func() {
		var actual MockModel
		err := suite.Model.Get(context.Background(), "mock", "t", &actual)
		require.EqualError(err, "redis: nil")
	})

	require.NoError(suite.DB.Del(context.Background(), "mock-test").Err())
}

func (suite *RedisModelHandlerSuite) TestSave() {
	require := suite.Require()

	expected := MockModel{Name: "save-test", Value: ""}

	require.NoError(suite.Model.Save(context.Background(), expected))

	v := suite.DB.Get(context.Background(), "mock-save-test").Val()

	var actual MockModel

	require.NoError(json.Unmarshal([]byte(v), &actual))
	require.Equal(expected, actual)
}

func (suite *RedisModelHandlerSuite) TestDelete() {
	require := suite.Require()

	expected := MockModel{Name: "test", Value: ""}
	v, _ := json.Marshal(expected)

	require.NoError(suite.DB.Set(context.Background(), "mock-test", v, 0).Err())

	suite.Run("testing successful delete", func() {
		require.NoError(suite.Model.Delete(context.Background(), "mock", "test"))

		err := suite.DB.Get(context.Background(), "mock-test").Err()
		require.EqualError(err, "redis: nil")
	})

	suite.Run("testing failed delete", func() {
		var actual MockModel

		err := suite.Model.Get(context.Background(), "mock", "t", &actual)
		require.EqualError(err, "redis: nil")
	})
}

func (suite *RedisModelHandlerSuite) TestUpdate() {
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

func TestRedisModelHandlerSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(RedisModelHandlerSuite))
}
