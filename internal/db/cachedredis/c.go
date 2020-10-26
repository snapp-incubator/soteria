package cachedredis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/internal/db/redis"
	"sync"
)

type ModelHandler struct {
	redisModelHandler *redis.ModelHandler
	cache             *cache.Cache
	validationTable   map[string]bool
	hit               int
	miss              int
	mu                sync.RWMutex
}

func NewCachedRedisModelHandler(redisModelHandler *redis.ModelHandler, cache *cache.Cache) *ModelHandler {
	return &ModelHandler{
		redisModelHandler: redisModelHandler,
		cache:             cache,
		validationTable:   make(map[string]bool),
		hit:               0,
		miss:              0,
		mu:                sync.RWMutex{},
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
	valid := mh.validationTable[key]
	item, found := mh.cache.Get(key)
	if found && valid {
		mh.mu.Lock()
		mh.hit++
		mh.mu.Unlock()
		json.Unmarshal([]byte(item.(string)), &v)
		return nil
	}

	err := mh.redisModelHandler.Get(ctx, modelName, pk, v)
	if err != nil {
		return err
	}
	value, _ := json.Marshal(v)
	mh.cache.SetDefault(key, string(value))

	mh.mu.Lock()
	mh.validationTable[key] = true
	mh.miss++
	mh.mu.Unlock()

	return nil
}

// Update invalidates cache and then finds and updates a model in redis
func (mh *ModelHandler) Update(ctx context.Context, model db.Model) error {
	key := redis.GenerateKey(model.GetMetadata().ModelName, model.GetPrimaryKey())

	mh.mu.Lock()
	mh.validationTable[key] = false
	mh.mu.Unlock()

	return mh.redisModelHandler.Update(ctx, model)
}

func (mh *ModelHandler) GetHit() int {
	mh.mu.Lock()
	hit := mh.hit
	mh.mu.Unlock()
	return hit
}

func (mh *ModelHandler) GetMiss() int {
	mh.mu.Lock()
	miss := mh.miss
	mh.mu.Unlock()
	return miss
}

func (mh *ModelHandler) IsCacheValid(key string) (bool, error) {
	mh.mu.Lock()
	value, ok := mh.validationTable[key]
	if !ok {
		return false, fmt.Errorf("key %v does not exist in cache validation table", key)
	}
	mh.mu.Unlock()
	return value, nil
}
