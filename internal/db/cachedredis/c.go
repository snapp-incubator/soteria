package cachedredis

import (
	"context"
	"encoding/json"
	"github.com/patrickmn/go-cache"
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/internal/db/redis"
)

type ModelHandler struct {
	redisModelHandler *redis.ModelHandler
	cache             *cache.Cache
}

func NewCachedRedisModelHandler(redisModelHandler *redis.ModelHandler, cache *cache.Cache) *ModelHandler {
	return &ModelHandler{
		redisModelHandler: redisModelHandler,
		cache:             cache,
	}
}

// Save saves a model in redis
func (mh *ModelHandler) Save(ctx context.Context, model db.Model) error {
	return mh.redisModelHandler.Save(ctx, model)
}

// Delete finds and deletes a model from redis and cache
func (mh *ModelHandler) Delete(ctx context.Context, modelName, pk string) error {
	key := redis.GenerateKey(modelName, pk)
	mh.cache.Delete(key)
	return mh.redisModelHandler.Delete(ctx, modelName, pk)
}

// Get searches cache at first then returns a model from redis or from cache, if exists
func (mh *ModelHandler) Get(ctx context.Context, modelName, pk string, v interface{}) error {
	key := redis.GenerateKey(modelName, pk)
	item, found := mh.cache.Get(key)
	if found {
		json.Unmarshal([]byte(item.(string)), &v)
		return nil
	}

	err := mh.redisModelHandler.Get(ctx, modelName, pk, v)
	if err != nil {
		return err
	}
	value, _ := json.Marshal(v)
	mh.cache.SetDefault(key, string(value))

	return nil
}

// Update invalidates cache and then finds and updates a model in redis
func (mh *ModelHandler) Update(ctx context.Context, model db.Model) error {
	key := redis.GenerateKey(model.GetMetadata().ModelName, model.GetPrimaryKey())

	mh.cache.Delete(key)

	return mh.redisModelHandler.Update(ctx, model)
}

