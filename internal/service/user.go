package service

import (
	"Vchat/internal/domain"
	"Vchat/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrUserNotFound          = repository.ErrUserNotFound
	ErrInvalidUserOrPassword = errors.New("用户不存在或密码错误")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Signup(ctx context.Context, domain domain.UserDomain) error {
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

func (s *UserService) Login(ctx context.Context, email string, password string) (*domain.UserDomain, error) {
	dmu, err := s.repo.FindByEmail(ctx, email)
	if err == ErrUserNotFound {
		return nil, repository.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(dmu.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidUserOrPassword
	}
	return dmu, nil
}

func (s *UserService) Profile(ctx context.Context, uid int64) (*domain.UserDomain, error) {
	dmu, err := s.repo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return dmu, nil
}

func (s *UserService) Edit(ctx context.Context, domain *domain.UserDomain) error {
	err := s.repo.Update(ctx, domain)
	return err
}
