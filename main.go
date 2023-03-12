package main

import (
	"time"

	"github.com/SonataUtopia/IM/models"
	"github.com/SonataUtopia/IM/router"
	"github.com/SonataUtopia/IM/utils"
	"github.com/spf13/viper"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	InitTimer()
	// test.TsetMySql()
	r := router.Router()
	r.Run(viper.GetString("port.server"))
}

func InitTimer() {
	utils.Timer(time.Duration(viper.GetInt("timeout.DelayHearbeat"))*time.Second, time.Duration(viper.GetInt("timeout.HeartbeatHz"))*time.Second, models.CleanConnection, "")
}
