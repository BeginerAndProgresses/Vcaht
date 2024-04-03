package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	// 自动注册，并不建议，以为可能不安全
	// 没有显式的SQL语句，建表也不知道建的是个是啥
	return db.AutoMigrate(&User{})
}
