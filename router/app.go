package router

import (
	"github.com/SonataUtopia/IM/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	//静态资源
	r.Static("/asset", "asset/")
	r.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	r.LoadHTMLGlob("views/**/*")

	//首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/ToRegister", service.ToRegister)
	r.GET("/ToChat", service.ToChat)
	r.GET("/Chat", service.Chat)
	r.POST("/LoadFriends", service.LoadFriends)

	//用户模块
	r.POST("/user/GetUserList", service.GetUserList)
	r.POST("/user/CreateUser", service.CreateUser)
	r.POST("/user/DeleteUser", service.DeleteUser)
	r.POST("/user/UpdateUser", service.UpdateUser)
	r.POST("/user/FindUserByNameAndPassword", service.FindUserByNameAndPassword)
	r.POST("/user/FindUserByID", service.FindUserByID)

	//发送消息
	r.GET("/user/SendMsg", service.SendMsg)
	r.GET("/user/SendUserMsg", service.SendUserMsg)
	r.POST("/user/redisMsg", service.RedisMsg)
	r.POST("/attach/Upload", service.Upload)

	//社交关系
	r.POST("/contact/Addfriend", service.AddFriend)
	r.POST("/contact/CreateCommunity", service.CreateCommunity)
	r.POST("/contact/LoadCommunity", service.LoadCommunity)
	r.POST("/contact/JoinGroup", service.JoinGroups)

	return r
}
