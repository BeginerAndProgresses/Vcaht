package service

import (
	"Vchat/internal/repository"
	"Vchat/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeVerifyTooMany

type CodeService interface {
	// Send 发送验证码，其中biz 是business 的缩写，一般代表业务，不知道什么业务就用这个代表
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type codeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo repository.CodeRepository, sms sms.Service) CodeService {
	return &codeService{
		repo: repo,
		sms:  sms,
	}
}

func (c *codeService) Send(ctx context.Context, biz, phone string) error {
	code := c.generate()
	err := c.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	const codeTplId = "123456"
	return c.sms.Send(ctx, codeTplId, []string{code}, phone)
}

func (c *codeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	ok, err := c.repo.Verify(ctx, biz, phone, inputCode)
	if err != nil {
		return false, err
	}
	return ok, err
}

func (c *codeService) generate() string {
	code := rand.Intn(100000)
	return fmt.Sprintf("%06d", code)
}
