package es

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Object_server_metadata struct {
	ID        uint           `gorm:"primaryKey"`
	Name string
	Version int64
	Size int64
	Hash string
}


// 连接数据库

func getEngine() *gorm.DB{
	dsn := "object_server:123.com@tcp(10.10.10.197:9000)/db_platform?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		//fmt.Println("connect to mysql is bad, err is", err)
		return nil
	}
	return db
}



