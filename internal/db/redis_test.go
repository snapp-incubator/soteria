package db

import (
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisModelHandler_Get(t *testing.T) {
	r := newTestRedis()
	s := RedisModelHandler{Client: r}
	expected := MockModel{Name: "test"}
	v, _ := json.Marshal(expected)

	err := r.Set("mock-test", v, 0).Err()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("testing successful get", func(t *testing.T) {
		var actual MockModel
		err := s.Get("mock", "test", &actual)
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})
	t.Run("testing failed get", func(t *testing.T) {
		var actual MockModel
		err := s.Get("mock", "t", &actual)
		assert.Error(t, err)
		assert.Equal(t, "redis: nil", err.Error())
	})

}

func TestRedisModelHandler_Save(t *testing.T) {
	r := newTestRedis()
	s := RedisModelHandler{Client: r}

	t.Run("testing save model", func(t *testing.T) {
		expected := MockModel{Name: "save-test"}
		err := s.Save(expected)
		assert.NoError(t, err)
		v := r.Get("mock-save-test").Val()
		var actual MockModel
		json.Unmarshal([]byte(v), &actual)
		assert.Equal(t, expected, actual)
	})
}

func TestRedisModelHandler_Delete(t *testing.T) {
	r := newTestRedis()
	s := RedisModelHandler{Client: r}
	expected := MockModel{Name: "test"}
	v, _ := json.Marshal(expected)

	err := r.Set("mock-test", v, 0).Err()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("testing successful delete", func(t *testing.T) {
		err := s.Delete("mock", "test")
		assert.NoError(t, err)
		err = r.Get("mock-test").Err()
		assert.Error(t, err)
		assert.Equal(t, "redis: nil", err.Error())
	})
	t.Run("testing failed delete", func(t *testing.T) {
		var actual MockModel
		err := s.Get("mock", "t", &actual)
		assert.Error(t, err)
		assert.Equal(t, "redis: nil", err.Error())
	})
}

<<<<<<< HEAD
=======
func TestRedisModelHandler_Update(t *testing.T) {
	r := newTestRedis()
	s := RedisModelHandler{Client: r}

	m := MockModel{Name: "test", Value: "test-1"}
	v, _ := json.Marshal(m)

	err := r.Set("mock-test", v, 0).Err()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("testing successful update", func(t *testing.T) {
		newModel := MockModel{Name: "test", Value: "test-2"}
		err = s.Update(newModel)
		assert.NoError(t, err)

		var updatedModel MockModel
		err = s.Get("mock", "test", &updatedModel)
		assert.NoError(t, err)

		assert.Equal(t, newModel, updatedModel)
	})
}

>>>>>>> master
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
<<<<<<< HEAD
	Name string `json:"name"`
=======
	Name  string `json:"name"`
	Value string `json:"value"`
>>>>>>> master
}

func (m MockModel) GetMetadata() MetaData {
	return MetaData{
		ModelName:    "mock",
		DateCreated:  time.Time{},
		DateModified: time.Time{},
	}
}

func (m MockModel) GetPrimaryKey() string {
	return m.Name
}
