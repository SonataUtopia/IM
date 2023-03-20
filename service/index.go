package service

import (
	"strconv"
	"text/template"

	"github.com/SonataUtopia/IM/models"
	"github.com/gin-gonic/gin"
)

func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("asset/html/login.html", "asset/html/head.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("asset/html/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("asset/html/index.html",
		"asset/html/head.html",
		"asset/html/foot.html",
		"asset/html/tabmenu.html",
		"asset/html/concat.html",
		"asset/html/group.html",
		"asset/html/profile.html",
		"asset/html/createcom.html",
		"asset/html/userinfo.html",
		"asset/html/main.html")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token

	ind.Execute(c.Writer, "chat")
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
