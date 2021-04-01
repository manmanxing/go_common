package util

import "github.com/jinzhu/gorm"

//封装一些对 gorm 的函数

func QueryForUpdate(db *gorm.DB) *gorm.DB {
	return db.Set("gorm:query_option", " FOR UPDATE ")
}
func QueryLockInShare(db *gorm.DB) *gorm.DB {
	return db.Set("gorm:query_option", " LOCK IN SHARE MODE ")
}