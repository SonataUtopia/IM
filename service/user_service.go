package service

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/SonataUtopia/IM/models"
	"github.com/SonataUtopia/IM/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"Code", "message"}
// @Router /user/getUserList [post]
func GetUserList(c *gin.Context) {
	// data := make([]*models.UserBasic, 10)
	data := models.GetUserList()

	models.GetUserList()
	c.JSON(200, gin.H{
		"Msg": data,
	})
}

// FindUserByNameAndPassword
// @Summary 查找用户
// @Tags 用户模块
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @Success 200 {string} json{"Code", "message"}
// @Router /user/FindUserByNameAndPassword [post]
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

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param repassword formData string false "确认密码"
// @Success 200 {string} json{"Code", "message"}
// @Router /user/CreateUser [post]
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

	// user.Password = password
	user.Password = utils.MakePassword(password, salt)
	user.Salt = salt
	// utils.DB.Create(user)
	models.CreateUser(user) //原案
	c.JSON(200, gin.H{
		"Code": 0,
		"Msg":  "密码正确",
		"Data": user,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id formData string false "用户ID"
// @Success 200 {string} json{"Code", "message"}
// @Router /user/DeleteUser [post]
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

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "用户ID"
// @param name formData string false "新用户名"
// @param password formData string false "新密码"
// @param phone formData string false "新手机号码"
// @param email formData string false "新邮箱"
// @Success 200 {string} json{"Code", "message"}
// @Router /user/UpdateUser [post]
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

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		fmt.Println(err)
	}
	tm := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println("MsgHandler error:", err)
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := models.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOKList(c.Writer, "ok", res)
}

func SearchFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriend(uint(id))
	// c.JSON(200, gin.H{
	// 	"Code":    0,
	// 	"Msg": "查询好友成功",
	// 	"Data":    users,
	// })
	utils.RespOKList(c.Writer, users, len(users))
}

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

func JoinGroups(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")

	//	name := c.Request.FormValue("name")
	data, msg := models.JoinGroup(uint(userId), comId)
	if data == 0 {
		utils.RespOK(c.Writer, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)

	} else {
		utils.RespFail(c.Writer, msg)
	}

}
