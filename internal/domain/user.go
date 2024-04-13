package domain

import "time"

type UserDomain struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	Nickname string
	Birthday time.Time
	AboutMe  string

	// 创建时间
	Ctime time.Time
	Utime time.Time
}
