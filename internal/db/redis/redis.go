package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
)

// RedisModelHandler implements ModelHandler interface
type ModelHandler struct {
	Client redis.Cmdable
}

// Save saves a model in redis
func (rmh ModelHandler) Save(ctx context.Context, model db.Model) error {
	md := model.GetMetadata()
	pk := model.GetPrimaryKey()
	key := GenerateKey(md.ModelName, pk)
	value, err := json.Marshal(model)
	if err != nil {
		return err
	}

	if err := rmh.Client.Set(ctx, key, string(value), 0).Err(); err != nil {
		return fmt.Errorf("%w: %s", db.ErrDb, err)
	}

	return nil
}

// Delete finds and deletes a model from redis and cache
func (rmh ModelHandler) Delete(ctx context.Context, modelName, pk string) error {
	key := GenerateKey(modelName, pk)

	res, err := rmh.Client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("%w: %s", db.ErrDb, err)
	}
	if res < 1 {
		return fmt.Errorf("key does not exist")
	}
	return nil
}

// Get returns a model from redis or from cache, if exists
func (rmh ModelHandler) Get(ctx context.Context, modelName, pk string, v interface{}) error {
	key := GenerateKey(modelName, pk)

	res, err := rmh.Client.Get(ctx, key).Result()
	if err != nil && err == redis.Nil {
		return fmt.Errorf("%s", err)
	} else if err != nil {
		return fmt.Errorf("%w: %s", db.ErrDb, err)
	}

	if err := json.Unmarshal([]byte(res), &v); err != nil {
		return err
	}

	return nil
}

// Update finds and updates a model in redis
func (rmh ModelHandler) Update(ctx context.Context, model db.Model) error {
	md := model.GetMetadata()
	pk := model.GetPrimaryKey()

	key := GenerateKey(md.ModelName, pk)

	value, err := json.Marshal(model)
	if err != nil {
		return err
	}

	pipeline := rmh.Client.Pipeline()
	pipeline.Del(ctx, key)
	pipeline.Set(ctx, key, string(value), 0)
	if _, err := pipeline.Exec(ctx); err != nil {
		return fmt.Errorf("%w: %s", db.ErrDb, err)
	}

	return nil
}

// GenerateKey is used to generate redis keys
func GenerateKey(modelName, pk string) string {
	return fmt.Sprintf("%v-%v", modelName, pk)
}
