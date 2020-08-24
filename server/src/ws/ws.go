package ws

import (
	"../Proto/Common"
	"fmt"
	"github.com/gin-gonic/gin" // web框架
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"sync"
)
/**
Manager所有websocket信息
 */
type Manager struct {
	Group map[string]map[string]*Client
	groupCount, clientCount uint
	Lock sync.Locker
	Register, UnRegister chan *Client
	Message chan *MessageData
	GroupMessage chan *GroupMessageData
	BroadCastMessage chan *BroadCastMessageData
}
/**
 单个websocket信息
 */
type Client struct {
	Id, Group string
	Socket *websocket.Conn
	Message chan []byte
}
/**
 单个发送数据
 */
type MessageData struct {
	Id, Group string
	Message []byte
}
/**
 组发送数据
 */
type GroupMessageData struct {
	Group string
	Message []byte
}
/**
 广播数据信息
 */
type BroadCastMessageData struct {
	Message []byte
}



/**
 获取manager
 */
var WebsocketManager = Manager{
	Group:            make(map[string]map[string]*Client),
	groupCount:       0,
	clientCount:      0,
	Register:         make(chan *Client, 128),
	UnRegister:       make(chan *Client, 128),
	Message:          make(chan *MessageData, 128),
	GroupMessage:     make(chan *GroupMessageData, 128),
	BroadCastMessage: make(chan *BroadCastMessageData, 128),
}

/**
 注册Client
 */
func (m *Manager) RegisterClient (client *Client)  {
	m.Register <- client
}
/**
 反注册Client
 */
func (m *Manager) UnRegisterClient (client *Client)  {
	m.UnRegister <- client
}

/**
 开始
 */
func (m *Manager) Start ()  {
	fmt.Println("websocket manager start")
	go m.SendService()
	go m.SendGroupService()
	go m.SendBroadService()
	for  {
		select {
		case client := <-m.Register: // 注册单个client
			//m.Lock.Lock()
			if m.Group[client.Group] == nil {
				m.Group[client.Group] = make(map[string]*Client)
				m.groupCount++
			}
			m.Group[client.Group][client.Id] = client
			m.clientCount++
			//m.Lock.Unlock()
		case client := <-m.UnRegister: // 注销单个client
			//m.Lock.Lock()
			if _, ok := m.Group[client.Group]; ok {
				if _, ok := m.Group[client.Group][client.Id]; ok {
					close(client.Message)
					delete(m.Group[client.Group], client.Id)
					m.clientCount--
					if len(m.Group[client.Group]) == 0 {
						delete(m.Group, client.Group)
						m.groupCount--
					}
				}
			}
			//m.Lock.Unlock()
		}
	}
}

/*
 读取数据
 */
func (c *Client) Read ()  {
	defer func() {
		WebsocketManager.UnRegister <- c
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] disconnect err: %s", c.Id, err)
		}
	}()
	for  {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		fmt.Println("数据：", message)

		// 以下为测试
		var newMsg []byte
		msgId := message[0]
		for index := 1; index < len(message); index++{
			newMsg = append(newMsg, message[index])
		}
		code := &msg.Code{}
		proto.Unmarshal(newMsg, code)
		fmt.Println("msgId: ", msgId)
		fmt.Println("msg: ", code.Msg)
		WebsocketManager.Send(c.Group, c.Id, code, msg.Event(msgId))
		//c.Message <- message
	}
}

/**
 写入数据
 */
func (c *Client) Write ()  {
	defer func() {
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] disconnect err: %s", c.Id, err)
		}
	}()
	for  {
		select {
		case message, ok := <- c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				break
			}
			//log.Printf("client [%s] write message: %s", c.Id, string(message))
			err := c.Socket.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Printf("client [%s] writemessage err: %s", c.Id, err)
			}
		}
	}
}

/**
 处理单个Client发送消息
 */
func (m *Manager) SendService ()  {
	for  {
		select {
		case msg := <- m.Message:
			if _, ok := m.Group[msg.Group]; ok {
				if _, ok := m.Group[msg.Group][msg.Id]; ok {
					m.Group[msg.Group][msg.Id].Message <- msg.Message
				}
			}
		}
	}
}

/**
 处理一个Group发送消息
 */
func (m *Manager) SendGroupService ()  {
	for  {
		select {
		case msg := <- m.GroupMessage:
			if _, ok := m.Group[msg.Group]; ok {
				for _, clientMsg := range m.Group[msg.Group]{
					clientMsg.Message <- msg.Message
				}
			}
		}
	}
}

/**
 发送广播消息
 */
func (m *Manager) SendBroadService ()  {
	for  {
		select {
		case msg := <- m.BroadCastMessage:
			if len(m.Group) > 0 {
				for _, group := range m.Group{
					for _, client := range group{
						client.Message <- msg.Message
					}
				}
			}
		}
	}
}

/**
 向指定client发送消息
 */
func (m *Manager) Send (group string, id string, message proto.Message, msgId msg.Event)  {
	sendMsg := messageChange(message, msgId)
	data := &MessageData{
		Id:      id,
		Group:   group,
		Message: sendMsg,
	}
	m.Message <- data
}

/**
 向指定Group发送消息
 */
func (m *Manager) SendGroup (group string, message proto.Message, msgId msg.Event)  {
	sendMsg := messageChange(message, msgId)
	data := &GroupMessageData{
		Group:   group,
		Message: sendMsg,
	}
	m.GroupMessage <- data
}

/**
 广播消息
 */
func (m *Manager) SendBroad (message proto.Message, msgId msg.Event)  {
	sendMsg := messageChange(message, msgId)
	data := &BroadCastMessageData{Message:sendMsg}
	m.BroadCastMessage <- data
}

/**
 协议加密成二进制
 */
func messageChange(msg proto.Message, msgId msg.Event) []byte {
	msgByte, _ := proto.Marshal(msg)
	var sendMsg []byte
	sendMsg = append(sendMsg, byte(msgId))
	sendMsg = append(sendMsg, msgByte...)
	return sendMsg
}

/**
 返回group数
 */
func (m *Manager) LenGroup () uint {
	return m.groupCount
}

/**
 返回client数
 */
func (m * Manager) LenClient () uint {
	return m.clientCount
}

/**
 监听事件（handler）
 */
func (m *Manager) WsClient (c *gin.Context)  {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket connect error: %s", c.Param("channel"))
		return
	}
	client := &Client{
		Id:      uuid.Must(uuid.NewV4()).String(),
		Group:   c.Param("channel"),
		Socket:  conn,
		Message: make(chan []byte),
	}
	m.RegisterClient(client)
	go client.Read()
	go client.Write()
}