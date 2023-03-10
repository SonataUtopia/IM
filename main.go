package main

import (
	"github.com/SonataUtopia/IM/router"
	"github.com/SonataUtopia/IM/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	// test.TsetMySql()
	r := router.Router()
	r.Run(":8080")
}
