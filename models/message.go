package models

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/SonataUtopia/IM/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// 消息
type Message struct {
	gorm.Model
	UserId     int64
	TargetId   int64
	Type       int    //消息类型			1.私聊	2.群聊	3.心跳
	Media      int    //消息内容类型		1.文字	2.表情包	3.图片	4.音频
	Content    string //消息内容
	CreateTime uint64 //创建时间
	ReadTime   uint64 //读取时间
	Pic        string
	Url        string
	Desc       string
	Amount     int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn          *websocket.Conn //连接
	Addr          string          //客户端地址
	FirstTime     uint64          //首次连接时间
	HeartbeatTime uint64          //心跳时间
	LoginTime     uint64          //登录时间
	DataQueue     chan []byte     //消息
}

var clientMap map[int64]*Node = make(map[int64]*Node, 0)

var rwLocker sync.RWMutex

// 聊天功能启动
func Chat(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	id := query.Get("userId")
	userId, _ := strconv.ParseInt(id, 10, 64)
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println("conn error:", err)
		return
	}

	//初始化node
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
	}

	//userId 与 node 绑定
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	//开启收发消息的协程
	go SendProc(node)
	go RecvProc(node)

}

// 发送消息
func SendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println("sendProc write error:", err)
				return
			}
		}
	}
}

// 接受消息
func RecvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println("recvProc read error:", err)
			return
		}

		msg := Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println("recvProc unmarshal error:", err)
		}

		if msg.Type == 3 {
			currentTime := uint64(time.Now().Unix())
			node.Heartbeat(currentTime)
		} else {
			Dispatch(data)
			BroadMsg(data)
		}

	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func BroadMsg(data []byte) {
	udpsendChan <- data
}

// 后端调度逻辑处理
func Dispatch(data []byte) {
	msg := Message{}
	msg.CreateTime = uint64(time.Now().Unix())
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1:
		SendMsg(msg.TargetId, data)
	case 2:
		SendGroupMsg(msg.TargetId, data)
	}
}

func SendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	targetIdStr := strconv.Itoa(int(userId))
	userIdStr := strconv.Itoa(int(jsonMsg.UserId))
	jsonMsg.CreateTime = uint64(time.Now().Unix())
	if ok {
		node.DataQueue <- msg
		node.Heartbeat(uint64(time.Now().Unix()))
	}

	//存储消息记录
	var key string
	if userId > jsonMsg.UserId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	res, err := utils.Red.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println("res error:", err)
	}
	score := float64(cap(res)) + 1
	_, e := utils.Red.ZAdd(ctx, key, &redis.Z{score, msg}).Result() //jsonMsg
	if e != nil {
		fmt.Println("ress error:", e)
	}
}

func SendGroupMsg(targetId int64, msg []byte) {
	userIds := SearchUserByGroupId(uint(targetId))
	for i := 0; i < len(userIds); i++ {
		//排除给自己的
		if targetId != int64(userIds[i]) {
			SendMsg(int64(userIds[i]), msg)
		}

	}
}

// 获取缓存里面的消息
func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))
	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}

	var rels []string
	var err error
	if isRev {
		rels, err = utils.Red.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = utils.Red.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println(err) //没有找到
	}
	return rels
}

// 更新用户心跳
func (node *Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
}

// 清理超时连接
func CleanConnection(param interface{}) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanConnection err", r)
		}
	}()
	currentTime := uint64(time.Now().Unix())
	for i := range clientMap {
		node := clientMap[i]
		if node.IsHeartbeatTimeOut(currentTime) {
			node.Conn.Close()
			return false
		}
	}
	return true
}

// 用户心跳是否超时
func (node *Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
	if node.HeartbeatTime+viper.GetUint64("timeout.HeartbeatMaxTime") <= currentTime {
		timeout = true
	}
	return
}
