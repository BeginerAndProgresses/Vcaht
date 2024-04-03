package repository

import (
	"Vchat/internal/domain"
	"Vchat/internal/repository/dao"
	"context"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateInsert
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{dao: dao}
}

func (r *UserRepository) Create(ctx context.Context, domain domain.UserDomain) error {
	return r.dao.Insert(ctx, dao.User{Email: domain.Email,
		Password: domain.Password})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain *domain.UserDomain, err error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return r.toDomain(u), nil
}

func (r *UserRepository) FindByUid(ctx context.Context, uid int64) (domain *domain.UserDomain, err error) {
	u, err := r.dao.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return r.toDomain(u), nil
}

func (r *UserRepository) toDomain(u *dao.User) *domain.UserDomain {
	return &domain.UserDomain{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
	}
}

func (r *UserRepository) Update(c context.Context, domain *domain.UserDomain) error {
	return r.dao.Update(c, domain)
}
