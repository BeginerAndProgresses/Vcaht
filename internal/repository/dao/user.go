package dao

import (
	"Vchat/internal/domain"
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (u User, err error)
	FindByUid(ctx context.Context, uid int64) (u User, err error)
	FindByPhone(ctx context.Context, phone string) (u User, err error)
	UpdateById(ctx context.Context, domain domain.UserDomain) error
}

type userDao struct {
	db *gorm.DB
}

var (
	ErrDuplicateInsert = errors.New("重复插入错误")
	ErrRecordNotFound  = gorm.ErrRecordNotFound
)

func NewUserDao(db *gorm.DB) UserDao {
	return &userDao{db: db}
}

func (d *userDao) Insert(ctx context.Context, u User) error {
	unixMilli := time.Now().UnixMilli()
	u.CAtime = unixMilli
	u.UAtime = unixMilli
	err := d.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateError uint16 = 1106
		if me.Number == duplicateError {
			return ErrDuplicateInsert
		}
	}
	return err
}

func (d *userDao) FindByEmail(ctx context.Context, email string) (u User, err error) {
	err = d.db.WithContext(ctx).Where("email = ?", email).Find(&u).Error
	return u, err
}

func (d *userDao) FindByUid(ctx context.Context, uid int64) (u User, err error) {
	err = d.db.WithContext(ctx).Where("id = ?", uid).Find(&u).Error
	return u, err
}

func (d *userDao) UpdateById(ctx context.Context, domain domain.UserDomain) error {
	user := d.toUser(domain)
	user.CAtime = time.Now().UnixMilli()
	err := d.db.WithContext(ctx).Model(User{}).Where("id", user.Id).Updates(&user).Error
	return err
}

func (d *userDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var res User
	err := d.db.WithContext(ctx).Where("phone", phone).First(&res).Error
	return res, err
}

func (d *userDao) toUser(domain domain.UserDomain) User {
	return User{
		Id:       domain.Id,
		Nickname: domain.Nickname,
		Birthday: domain.Birthday.UnixMilli(),
		AboutMe:  domain.AboutMe,
		CAtime:   domain.Ctime.UnixMilli(),
		UAtime:   domain.Utime.UnixMilli(),
		Phone:    sql.NullString{String: domain.Phone, Valid: domain.Phone != ""},
		Email:    sql.NullString{String: domain.Email, Valid: domain.Email != ""},
	}
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `json:"email" gorm:"unique"` // 设置唯一索引
	Password string         `json:"password"`
	Phone    sql.NullString `json:"phone" gorm:"unique"`
	Nickname string         `json:"nickname"`
	Birthday int64          `json:"birthday"`
	AboutMe  string         `json:"aboutMe"`
	CAtime   int64
	UAtime   int64
}
