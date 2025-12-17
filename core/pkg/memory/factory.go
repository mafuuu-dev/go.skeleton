package memory

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Factory struct {
	ctx   context.Context
	cache *redis.Client
}

func New(ctx context.Context, cache *redis.Client) *Factory {
	return &Factory{
		ctx:   ctx,
		cache: cache,
	}
}

func (f *Factory) Context() context.Context {
	return f.ctx
}

func (f *Factory) Cache() *redis.Client {
	return f.cache
}
