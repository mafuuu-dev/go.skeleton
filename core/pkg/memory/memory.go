package memory

import (
	"backend/core/pkg/errorsx"
	"backend/core/types"
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Memory[T any] struct {
	ctx context.Context
	db  *redis.Client
	key string
}

func NewMemory[T any](ctx context.Context, db *redis.Client, key types.StoreKey, args []string) *Memory[T] {
	postfix := ""
	for _, arg := range args {
		postfix += ":" + arg
	}

	return &Memory[T]{ctx: ctx, db: db, key: string(key) + postfix}
}

func (memory *Memory[T]) Get() (*T, error) {
	data, err := memory.db.Get(memory.ctx, memory.key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, errorsx.Error(err)
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, errorsx.Error(err)
	}

	return &result, nil
}

func (memory *Memory[T]) Set(value *T, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return errorsx.Error(err)
	}

	return memory.db.Set(memory.ctx, memory.key, bytes, expiration).Err()
}

func (memory *Memory[T]) Delete() error {
	return memory.db.Del(memory.ctx, memory.key).Err()
}

func (memory *Memory[T]) Exists() (bool, error) {
	count, err := memory.db.Exists(memory.ctx, memory.key).Result()
	if err != nil {
		return false, errorsx.Error(err)
	}

	return count > 0, nil
}

func (memory *Memory[T]) ExpiresIn(ttl time.Duration) error {
	ok, err := memory.db.Expire(memory.ctx, memory.key, ttl).Result()
	if err != nil {
		return errorsx.Error(err)
	}
	if !ok {
		return redis.Nil
	}

	return nil
}
