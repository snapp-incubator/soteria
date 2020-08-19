package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
)

// RedisModelHandler implements ModelHandler interface
type RedisModelHandler struct {
	Client redis.Cmdable
}

// Save saves a model in redis
func (rmh RedisModelHandler) Save(model db.Model) error {
	md := model.GetMetadata()
	pk := model.GetPrimaryKey()
	key := generateKey(md.ModelName, pk)
	value, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return rmh.Client.Set(key, string(value), 0).Err()
}

// Save finds and deletes a model from redis
func (rmh RedisModelHandler) Delete(modelName, pk string) error {
	key := generateKey(modelName, pk)
	res, err := rmh.Client.Del(key).Result()
	if err != nil {
		return err
	}
	if res < 1 {
		return fmt.Errorf("key does not exist")
	}
	return nil
}

// Save finds and returns a model from redis
func (rmh RedisModelHandler) Get(modelName, pk string, v interface{}) error {
	key := generateKey(modelName, pk)
	res, err := rmh.Client.Get(key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(res), &v)
	if err != nil {
		return err
	}
	return nil
}

func (rmh RedisModelHandler) Update(model db.Model) error {
	md := model.GetMetadata()
	pk := model.GetPrimaryKey()

	key := generateKey(md.ModelName, pk)

	value, err := json.Marshal(model)
	if err != nil {
		return err
	}

	pipeline := rmh.Client.Pipeline()
	pipeline.Del(key)
	pipeline.Set(key, string(value), 0)
	_, err = pipeline.Exec()

	return err
}

func generateKey(modelName, pk string) string {
	return fmt.Sprintf("%v-%v", modelName, pk)
}
