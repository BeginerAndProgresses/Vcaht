package repository

import (
	"Vchat/internal/domain"
	"Vchat/internal/repository/cache"
	"Vchat/internal/repository/dao"
	"context"
	"database/sql"
	"time"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateInsert
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, domain domain.UserDomain) error
	FindByEmail(ctx context.Context, email string) (domain.UserDomain, error)
	FindByUid(ctx context.Context, uid int64) (domain.UserDomain, error)
	FindByPhone(ctx context.Context, phone string) (domain.UserDomain, error)
	Update(c context.Context, domain domain.UserDomain) error
}

type userRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &userRepository{dao: dao, cache: cache}
}

func (r *userRepository) Create(ctx context.Context, domain domain.UserDomain) error {
	return r.dao.Insert(ctx, r.toEntity(domain))
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (domain.UserDomain, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.UserDomain{}, err
	}
	return r.toDomain(u), nil
}

// FindByUid 缓存击穿保留该业务，可能这会导致数据库崩溃，导致其他业务无法执行
//func (r *UserRepository) FindByUid(ctx context.Context, uid int64) (dm *domain.UserDomain, err error) {
//	u, err := r.cache.Get(ctx, uid)
//	switch err {
//	case nil:
//		return u, err
//	case cache.ErrKeyNotExist:
//		ur, err := r.dao.FindByUid(ctx, uid)
//		if err != nil {
//			return &domain.UserDomain{}, err
//		}
//		du := r.toDomain(ur)
//		_ = r.cache.Set(ctx, *du)
//		return du, nil
//	}
//	return &domain.UserDomain{}, err
//}

// FindByUid 缓存击穿返回错误
func (r *userRepository) FindByUid(ctx context.Context, uid int64) (domain.UserDomain, error) {
	u, err := r.cache.Get(ctx, uid)
	if err == nil {
		return u, nil
	}
	ur, err := r.dao.FindByUid(ctx, uid)
	if err != nil {
		return domain.UserDomain{}, err
	}
	du := r.toDomain(ur)
	err = r.cache.Set(ctx, du)
	if err != nil {
		return domain.UserDomain{}, err
	}
	return du, nil
}

func (r *userRepository) toDomain(u dao.User) domain.UserDomain {
	return domain.UserDomain{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
	}
}

func (r *userRepository) toEntity(du domain.UserDomain) dao.User {
	return dao.User{
		Id:       du.Id,
		Email:    sql.NullString{String: du.Email, Valid: du.Email != ""},
		Phone:    sql.NullString{String: du.Phone, Valid: du.Phone != ""},
		Password: du.Password,
		Nickname: du.Nickname,
		Birthday: du.Birthday.UnixMilli(),
		AboutMe:  du.AboutMe,
	}
}

func (r *userRepository) Update(c context.Context, domain domain.UserDomain) error {
	return r.dao.UpdateById(c, domain)
}

func (r *userRepository) FindByPhone(c context.Context, phone string) (domain.UserDomain, error) {
	u, err := r.dao.FindByPhone(c, phone)
	if err != nil {
		return domain.UserDomain{}, err
	}
	return r.toDomain(u), nil
}
