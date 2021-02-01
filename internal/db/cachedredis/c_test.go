package cachedredis

import (
	"context"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	redisdb "gitlab.snapp.ir/dispatching/soteria/v3/internal/db/redis"
	"testing"
	"time"
)

func TestModelHandler_Get(t *testing.T) {
	r := newTestRedis()
	s := ModelHandler{
		redisModelHandler: &redisdb.ModelHandler{Client: r},
		cache:             cache.New(time.Minute*60, time.Minute*60),
	}
	expected := MockModel{Name: "test"}
	v, _ := json.Marshal(expected)

	err := r.Set(context.Background(), "mock-test", v, 0).Err()
	assert.NoError(t, err)

	t.Run("testing first successful get", func(t *testing.T) {
		var actual MockModel
		err := s.Get(context.Background(), "mock", "test", &actual)
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)

		s.cache.Get(redisdb.GenerateKey("mock", "test"))
	})
	t.Run("testing second successful get", func(t *testing.T) {
		var actual MockModel
		err := s.Get(context.Background(), "mock", "test", &actual)
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})
	t.Run("testing failed get", func(t *testing.T) {
		var actual MockModel
		err := s.Get(context.Background(), "mock", "t", &actual)
		assert.Error(t, err)
		assert.Equal(t, "redis: nil", err.Error())
	})
	assert.NoError(t, err)
	err = r.Del(context.Background(), "mock-test").Err()
	assert.NoError(t, err)
}

func TestModelHandler_Save(t *testing.T) {
	r := newTestRedis()
	s := ModelHandler{
		redisModelHandler: &redisdb.ModelHandler{Client: r},
		cache:             cache.New(time.Minute*60, time.Minute*60),
	}

	t.Run("testing save model", func(t *testing.T) {
		expected := MockModel{Name: "save-test"}
		err := s.Save(context.Background(), expected)
		assert.NoError(t, err)
		v := r.Get(context.Background(), "mock-save-test").Val()
		var actual MockModel
		json.Unmarshal([]byte(v), &actual)
		assert.Equal(t, expected, actual)
	})
}

func TestRedisModelHandler_Delete(t *testing.T) {
	r := newTestRedis()
	s := ModelHandler{
		redisModelHandler: &redisdb.ModelHandler{Client: r},
		cache:             cache.New(time.Minute*60, time.Minute*60),
	}

	expected := MockModel{Name: "test"}
	v, _ := json.Marshal(expected)

	err := r.Set(context.Background(), "mock-test", v, 0).Err()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("testing successful delete", func(t *testing.T) {
		err := s.Delete(context.Background(), "mock", "test")
		assert.NoError(t, err)
		err = r.Get(context.Background(), "mock-test").Err()
		assert.Error(t, err)
		assert.Equal(t, "redis: nil", err.Error())
		item, ok := s.cache.Get("mock-test")
		assert.False(t, ok)
		assert.Nil(t, item)

	})
	t.Run("testing failed delete", func(t *testing.T) {
		var actual MockModel
		err := s.Get(context.Background(), "mock", "t", &actual)
		assert.Error(t, err)
		assert.Equal(t, "redis: nil", err.Error())
	})
}

func TestRedisModelHandler_Update(t *testing.T) {
	r := newTestRedis()
	s := ModelHandler{
		redisModelHandler: &redisdb.ModelHandler{Client: r},
		cache:             cache.New(time.Minute*60, time.Minute*60),
	}

	m := MockModel{Name: "test", Value: "test-1"}
	v, _ := json.Marshal(m)

	err := r.Set(context.Background(), "mock-test", v, 0).Err()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("testing successful update", func(t *testing.T) {
		newModel := MockModel{Name: "test", Value: "test-2"}
		err = s.Update(context.Background(), newModel)
		assert.NoError(t, err)
		assert.NoError(t, err)
		var updatedModel MockModel
		err = s.Get(context.Background(), "mock", "test", &updatedModel)
		assert.NoError(t, err)
		assert.Equal(t, newModel, updatedModel)
	})
}

func newTestRedis() *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return client
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
