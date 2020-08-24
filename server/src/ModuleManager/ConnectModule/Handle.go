package ConnectModule

import (
	"../../Proto/Common"
	"../DbModule"
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
)

type handle struct {
	ws *websocket.Conn // ws连接
	wsMsg websocket.Codec
	user *DbModule.UserInfo
	infoDb *DbModule.InfoDb
}
const idAdd = 1000000

/**
 协议处理函数
 */
func (h *handle) handleFunc(msgId msg.Event, msgByte []byte)  {
	fmt.Println("收到的协议id: ", msgId)
	switch msgId {

	}
}

/**
 解析数据
 */
func (h *handle) decodeData (byteData []byte, receiveData proto.Message)  {
	proto.Unmarshal(byteData, receiveData)
}

/**
 获取code
 */
func (h *handle) getCode (codeType msg.CodeType, str string) *msg.Code {
	code := &msg.Code{}
	code.Msg = str
	code.Code = codeType
	return code
}

/**
 发送消息
 */
func (h *handle) send (msgId msg.Event, msgData proto.Message)  {
	var sendByte []byte
	msgDataByte, _ := proto.Marshal(msgData)
	sendByte = append(sendByte, byte(msgId))
	sendByte = append(sendByte, msgDataByte...)
	if err := websocket.Message.Send(h.ws, sendByte); err != nil {
		fmt.Println("发送出错", err.Error())
	}
}

/**
 id转换
 */
func changeId(id int) int {
	if id > idAdd {
		id = id - idAdd
	} else {
		id = id + idAdd
	}
	return id
}