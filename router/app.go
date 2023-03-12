package router

import (
	"github.com/SonataUtopia/IM/docs"
	"github.com/SonataUtopia/IM/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	//swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")

	//首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/ToRegister", service.ToRegister)
	r.GET("/ToChat", service.ToChat)
	r.GET("/Chat", service.Chat)
	r.POST("/SearchFriends", service.SearchFriends)

	//用户模块
	r.POST("/user/GetUserList", service.GetUserList)
	r.POST("/user/CreateUser", service.CreateUser)
	r.POST("/user/DeleteUser", service.DeleteUser)
	r.POST("/user/UpdateUser", service.UpdateUser)
	r.POST("/user/FindUserByNameAndPassword", service.FindUserByNameAndPassword)

	//发送消息
	r.GET("/user/SendMsg", service.SendMsg)
	r.GET("/user/SendUserMsg", service.SendUserMsg)
	r.POST("/attach/Upload", service.Upload)

	//社交关系
	r.POST("/contact/Addfriend", service.AddFriend)
	r.POST("/contact/CreateCommunity", service.CreateCommunity)
	r.POST("/contact/LoadCommunity", service.LoadCommunity)
	r.POST("/contact/JoinGroup", service.JoinGroups)

	return r
}
