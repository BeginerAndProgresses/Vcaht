package service

import (
	"Vchat/internal/domain"
	"Vchat/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUser
	ErrUserNotFound          = repository.ErrUserNotFound
	ErrInvalidUserOrPassword = errors.New("用户不存在或密码错误")
)

type UserService interface {
	FindOrCreate(ctx context.Context, phone string) (domain.UserDomain, error)
	Signup(ctx context.Context, domain domain.UserDomain) error
	Login(ctx context.Context, email string, password string) (domain.UserDomain, error)
	Profile(ctx context.Context, uid int64) (domain.UserDomain, error)
	Edit(ctx context.Context, domain domain.UserDomain) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) FindOrCreate(ctx context.Context, phone string) (domain.UserDomain, error) {
	// 查询，假设大部分用户已经存在
	u, err := s.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// 两种情况
		// err == nil, u是可用的
		// err != nil, 系统错误
		return u, err
	}
	// 用户不存在，创建
	err = s.repo.Create(ctx, domain.UserDomain{
		Phone: phone,
	})
	// 有两种情况，一是唯一冲突，二是系统错误
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.UserDomain{}, err
	}
	return s.repo.FindByPhone(ctx, phone)
}

func (s *userService) Signup(ctx context.Context, domain domain.UserDomain) error {
	//	加密：
	password, err := bcrypt.GenerateFromPassword([]byte(domain.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	domain.Password = string(password)
	err = s.repo.Create(ctx, domain)
	if err == ErrDuplicateEmail {
		return ErrDuplicateEmail
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) Login(ctx context.Context, email string, password string) (domain.UserDomain, error) {
	dmu, err := s.repo.FindByEmail(ctx, email)
	if err == ErrUserNotFound {
		return domain.UserDomain{}, repository.ErrUserNotFound
	}
	if err != nil {
		return domain.UserDomain{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(dmu.Password), []byte(password))
	if err != nil {
		return domain.UserDomain{}, ErrInvalidUserOrPassword
	}
	return dmu, nil
}

func (s *userService) Profile(ctx context.Context, uid int64) (domain.UserDomain, error) {
	dmu, err := s.repo.FindByUid(ctx, uid)
	if err != nil {
		return domain.UserDomain{}, err
	}
	return dmu, nil
}

func (s *userService) Edit(ctx context.Context, domain domain.UserDomain) error {
	err := s.repo.Update(ctx, domain)
	return err
}
