package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode        string
	ErrCodeSendTooMany   = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooMany = errors.New("验证次数过多")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type redisCodeCache struct {
	cmd redis.Cmdable
}

func (r *redisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := r.cmd.Eval(ctx, luaSetCode, []string{r.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case -2:
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		return ErrCodeSendTooMany
	default:
		return nil
	}
}

func (r *redisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := r.cmd.Eval(ctx, luaVerifyCode, []string{r.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case -2:
		return false, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	default:
		return true, nil
	}
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &redisCodeCache{cmd: cmd}
}

func (r *redisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

type localCodeCache struct {
	mc *MapCache
}

func NewLocalCodeCache(mc *MapCache) CodeCache {
	return &localCodeCache{mc: mc}
}

func (l *localCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	err := l.mc.Set(l.key(biz, phone), code)
	l.mc.SetExAt(time.Now().Add(time.Minute * 15))
	if err != nil {
		return ErrMyCacheIsNil
	}
	return nil
}

func (l *localCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	get, err := l.mc.Get(l.key(biz, phone))
	if err != nil {
		return false, ErrMyCacheIsNil
	}
	code := get.(string)
	if code != inputCode {
		return false, ErrCodeVerifyTooMany
	}
	return true, nil
}

func (l *localCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
