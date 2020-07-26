package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
)

type ModelHandler struct {
	Client redis.Cmdable
}

func (rmh ModelHandler) Save(model db.Model) error {
	md := model.GetMetadata()
	pk := model.GetPrimaryKey()
	key := generateKey(md.ModelName, pk)
	value, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return rmh.Client.Set(key, string(value), 0).Err()
}

func (rmh ModelHandler) Delete(modelName, pk string) error {
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

func (rmh ModelHandler) Get(modelName, pk string, v interface{}) error {
	key := generateKey(modelName, pk)
	res, err := rmh.Client.Get(key).Result()
	if err != nil {
		return nil
	}
	err = json.Unmarshal([]byte(res), &v)
	if err != nil {
		return err
	}
	return nil
}

func generateKey(modelName, pk string) string {
	return fmt.Sprintf("%v-%v", modelName, pk)
}
