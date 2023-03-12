package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/SonataUtopia/IM/utils"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

// 消息
type Message struct {
	gorm.Model
	FormId   int64
	TargetId int64
	Type     int    //消息类型			1.私聊	2.群聊	3.广播
	Media    int    //消息内容类型		1.文字	2.表情包	3.图片	4.音频
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

var clientMap map[int64]*Node = make(map[int64]*Node, 0)

var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	//检验token
	// token := query.Get("token")

	query := request.URL.Query()
	id := query.Get("userId")
	userId, _ := strconv.ParseInt(id, 10, 64)
	// fmt.Println("Id,userId-------------------------------------------------------", id, userId)
	// msgType := query.Get("msgType")
	// targetId := query.Get("targetId")
	// context := query.Get("context")
	isvalida := true //check token
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//获取token
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	//判断用户关系

	//userId 与 node 绑定
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	//发送消息
	go SendProc(node)

	//接收消息
	go RecvProc(node)

	SendMsg(userId, []byte("欢迎来到聊天室"))
}

func SendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws]sendProc >>>> msg :", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func RecvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		Dispatch(data)
		BroadMsg(data)
		// fmt.Println("[ws] <<<<<", data)
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func BroadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
	fmt.Println("init goroutine ")
}

// 完成udp发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			// fmt.Println("udpSendProc  data :", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成数据接收协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		// fmt.Println("udpRecvProc  data :", string(buf[0:n]))
		Dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func Dispatch(data []byte) {
	// fmt.Println("Dispatch onnnnnnnnnnn!!!")
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		SendMsg(msg.TargetId, data)
		// case 2://群发
		// 	sendGroupMsg()
		// case 3://广播
		// 	sendAllMsg()
		// case 4:

	}
}

func SendMsg(userId int64, msg []byte) {
	// fmt.Println("sendMsg >>> userId:", userId, "msg:", string(msg))
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	fmt.Println("node:", node, "\t\t\tok:", ok)
	if ok {
		// fmt.Println("SendMsg >>> userID: ", userId, "  msg:", string(msg))
		node.DataQueue <- msg
	}
}

func JoinGroup(userId uint, comId string) (int, string) {
	contact := Contact{}
	contact.OwnerId = userId
	//contact.TargetId = comId
	contact.Type = 2
	community := Community{}

	utils.DB.Where("id=? or name=?", comId, comId).Find(&community)
	if community.Name == "" {
		return -1, "没有找到群"
	}
	utils.DB.Where("owner_id=? and target_id=? and type =2 ", userId, comId).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1, "已加过此群"
	} else {
		contact.TargetId = community.ID
		utils.DB.Create(&contact)
		return 0, "加群成功"
	}
}
