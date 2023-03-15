package main

import (
	"fmt"

	"github.com/SonataUtopia/IM/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		fmt.Println("err:", err)
	}

	// 创建表
	db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.GroupBasic{})
	db.AutoMigrate(&models.Contact{})
	db.AutoMigrate(&models.Community{})
}
