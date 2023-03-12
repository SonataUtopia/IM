package main

import (
	"fmt"

	"github.com/SonataUtopia/IM/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// type Product struct {
// 	gorm.Model
// 	Name          string
// 	Password      string
// 	Phone         string
// 	Email         string
// 	Identity      string
// 	ClentIp       string
// 	ClientPort    string
// 	LoginTime     uint64
// 	HeartBeatTime uint64
// 	LogoutTime    uint64
// 	IsLogout      bool
// 	DeviceInfo    string
// }

func main() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		fmt.Println("err:", err)
	}

	// Migrate the schema
	// utils.DB.AutoMigrate(&models.UserBasic{})
	// db.AutoMigrate(&models.Message{})
	// db.AutoMigrate(&models.GroupBasic{})
	// db.AutoMigrate(&models.Contact{})
	db.AutoMigrate(&models.Community{})

	// Create
	// user := &models.UserBasic{}
	// user.Name = "createByUntils"
	// db.Create(user)

	// Read
	// fmt.Println(db.First(user, 1))
	// db.First(user, 1)                 // find product with integer primary key
	// db.First(user, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	// db.Model(user).Update("password", "123")
	// Update - update multiple fields
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	// db.Delete(&product, 1)
}
