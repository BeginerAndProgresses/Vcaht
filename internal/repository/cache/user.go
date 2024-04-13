package cache

import (
	"Vchat/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.UserDomain, error)
	Set(ctx context.Context, du domain.UserDomain) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: 30 * time.Minute,
	}
}

func (c *RedisUserCache) Get(ctx context.Context, uid int64) (domain.UserDomain, error) {
	key := c.key(uid)
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.UserDomain{}, err
	}
	//	使用json反序列化
	var userDomain domain.UserDomain
	err = json.Unmarshal([]byte(data), &userDomain)
	return userDomain, err
}

func (c *RedisUserCache) Set(ctx context.Context, du domain.UserDomain) error {
	key := c.key(du.Id)
	//	使用json序列化
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisUserCache) key(uid int64) string {
	return fmt.Sprintf("user-info-%d", uid)
}

type UserCacheV1 struct {
	client *redis.Client
}
