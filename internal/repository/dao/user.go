package dao

import (
	"Vchat/internal/domain"
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	db *gorm.DB
}

var (
	ErrDuplicateInsert = errors.New("重复插入错误")
	ErrRecordNotFound  = gorm.ErrRecordNotFound
)

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (d *UserDao) Insert(ctx context.Context, u User) error {
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

func (d *UserDao) FindByEmail(ctx context.Context, email string) (u *User, err error) {
	err = d.db.WithContext(ctx).Where("email = ?", email).Find(&u).Error
	return u, err
}

func (d *UserDao) FindByUid(ctx context.Context, uid int64) (u *User, err error) {
	err = d.db.WithContext(ctx).Where("id = ?", uid).Find(&u).Error
	return u, err
}

func (d *UserDao) Update(ctx context.Context, domain *domain.UserDomain) error {
	user := d.toUser(domain)
	user.CAtime = time.Now().UnixMilli()
	err := d.db.WithContext(ctx).Model(User{}).Where("id", user.Id).Updates(&user).Error
	return err
}

func (d *UserDao) toUser(domain *domain.UserDomain) *User {
	return &User{
		Id: domain.Id,
		//Password:
		Nickname: domain.Nickname,
		Birthday: domain.Birthday,
		AboutMe:  domain.AboutMe,
	}
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `json:"email" gorm:"unique"` // 设置唯一索引
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Birthday string `json:"birthday"`
	AboutMe  string `json:"aboutMe"`
	CAtime   int64
	UAtime   int64
}
