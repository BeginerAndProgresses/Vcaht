package repository

import (
	"Vchat/internal/repository/cache"
	"context"
)

var ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type cacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCacheCodeRepository(cache cache.CodeCache) CodeRepository {
	return &cacheCodeRepository{
		cache: cache,
	}
}
func (c *cacheCodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *cacheCodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, inputCode)
}
