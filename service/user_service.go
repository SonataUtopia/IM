package service

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/SonataUtopia/IM/models"
	"github.com/SonataUtopia/IM/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// 获取用户列表
func GetUserList(c *gin.Context) {
	data := models.GetUserList()

	models.GetUserList()
	c.JSON(200, gin.H{
		"Msg": data,
	})
}

// 通过用户名与密码寻找用户
func FindUserByNameAndPassword(c *gin.Context) {
	// data := models.UserBasic{}
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"Code": -1,
			"Msg":  "用户不存在",
			"Data": user,
		})
		return
	}

	flag := utils.ValidPassword(password, user.Salt, user.Password)
	if !flag {
		c.JSON(200, gin.H{
			"Code": -1,
			"Msg":  "密码不正确",
			"Data": user,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data := models.FindUserByNameAndPassword(name, pwd)

	c.JSON(200, gin.H{
		"Code": 0,
		"Msg":  "登录成功",
		"Data": data,
	})
}

// 通过ID寻找用户
func FindUserByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))

	data := models.FindUserByID(uint(userId))
	utils.RespOK(c.Writer, data, "ok")
}

func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(200, gin.H{
			"Code": -1,
			"Msg":  "用户名、密码不能为空",
			"Data": user,
		})
		return
	}
	if password != repassword {
		// fmt.Println(password, "|", repassword)
		c.JSON(200, gin.H{
			"Code": -1,
			"Msg":  "两次密码不一致",
			"Data": user,
		})
		return
	}

	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(200, gin.H{
			"Code": -1,
			"Msg":  "用户名已被使用",
			"Data": data,
		})
		return
	}

	user.Password = utils.MakePassword(password, salt)
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"Code": 0,
		"Msg":  "密码正确",
		"Data": user,
	})
}

// 删除用户
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Request.FormValue("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"Code": 0,
		"Msg":  "删除用户成功",
		"Data": user,
	})
}

// 更新用户信息
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Avatar = c.PostForm("icon")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		c.JSON(200, gin.H{
			"Code": -1,
			"Msg":  "格式不匹配",
			"Data": user,
		})
		return
	}

	models.UpdateUser(user)
	c.JSON(200, gin.H{
		"Code": 0,
		"Msg":  "修改用户成功",
		"Data": user,
	})
}

// 加载消息记录
func GetMsgLogging(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))
	targetId, _ := strconv.Atoi(c.PostForm("targetId"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	isCom, _ := strconv.ParseBool(c.PostForm("isCom"))
	res := models.GetMsgLogging(int64(userId), int64(targetId), int64(start), int64(end), isCom, isRev)
	utils.RespOKList(c.Writer, "ok", res)
}

// 加载好友列表
func LoadFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.LoadFriend(uint(id))
	utils.RespOKList(c.Writer, users, len(users))
}

// 添加好友
func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetName := c.Request.FormValue("targetName")
	Code, msg := models.AddFriend(uint(userId), targetName)
	if Code == 0 {
		utils.RespOK(c.Writer, Code, msg)

	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 创建群聊
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	icon := c.Request.FormValue("icon")
	desc := c.Request.FormValue("memo")
	community := models.Community{
		OwnerId: uint(ownerId),
		Name:    name,
		Icon:    icon,
		Desc:    desc,
	}

	Code, msg := models.CreateCommunity(community)
	if Code == 0 {
		utils.RespOK(c.Writer, Code, msg)

	} else {
		utils.RespFail(c.Writer, msg)
	}

}

// 加入群聊
func JoinGroups(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")

	data, msg := models.JoinGroup(uint(userId), comId)
	if data == 0 {
		utils.RespOK(c.Writer, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 加载群聊列表
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)

	} else {
		utils.RespFail(c.Writer, msg)
	}
}
