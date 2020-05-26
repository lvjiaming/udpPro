package ConnectModule

import (
	"../../Cfg"
	"../../ModuleManager"
	"../../Proto/Common"
	"../../Proto/JiZhangBo"
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
	case msg.Event_EVENT_REGISTER_REQ: // 注册
		h.dealWithRegister(msgByte)
	case msg.Event_EVENT_LOGIN_REQ: // 登录
		h.DealWithLogin(msgByte)
	case msg.Event_EVENT_ADD_ONE_INFO: // 新增信息
		h.addInfo(msgByte)
	case msg.Event_EVENT_DEL_INFO_REQ: // 删除信息
		h.delInfo(msgByte)
	case msg.Event_EVENT_QUERY_INFO_REQ: // 查询信息
		h.queryInfo(msgByte)
	case msg.Event_EVENT_CHANGE_INFO_REQ: // 修改信息
		h.changeInfo(msgByte)
	}
}

/**
 处理注册请求
 */
func (h *handle) dealWithRegister (msgByte []byte)  {
	reg := &jiZhangBo.RegisterReq{}
	h.decodeData(msgByte, reg)
	_, err := UserDb.QueryUserToTerm(DbModule.UserName(reg.User.Name))
	if err != nil {
		fmt.Print("玩家还未被注册，允许注册\n")
		err := UserDb.Insert(DbModule.UserName(reg.User.Name), DbModule.UserPwd(reg.User.Pwd))
		if err != nil {
			h.registerRep(msg.CodeType_ERR, "插入数据错误")
		} else {
			h.registerRep(msg.CodeType_SUC, "注册成功")
		}
	} else {
		fmt.Print("玩家已存在\n")
		h.registerRep(msg.CodeType_ERR, "玩家已存在")
	}
}

/**
 注册回复
 */
func (h *handle) registerRep (codeType msg.CodeType, str string)  {
	registerRep := &jiZhangBo.RegisterRep{
		Code: h.getCode(codeType, str),
	}
	h.send(msg.Event_EVENT_REGISTER_REP, registerRep)
}

/**
 处理登录请求
 */
func (h *handle) DealWithLogin (msgByte []byte)  {
	login := &jiZhangBo.LoginReq{}
	h.decodeData(msgByte, login)
	userInfo, err := UserDb.QueryUserToTerm(DbModule.UserName(login.User.Name))
	if err != nil {
		fmt.Println(err.Error())
		h.loginRep(msg.CodeType_ERR, "玩家不存在", nil)
	} else {
		if userInfo.Password != login.User.Pwd {
			h.loginRep(msg.CodeType_ERR, "密码不正确", nil)
		} else {
			hasLogin := false // 用于判断是否已经登陆
			for _, user := range userList{
				if user.user != nil && user.user.UserId == userInfo.UserId {
					hasLogin = true
				}
			}
			if hasLogin == true {
				fmt.Print("玩家已登录，重复登录")
				h.loginRep(msg.CodeType_ERR, "玩家已登录，重复登录", nil)
			} else {
				h.user = &userInfo
				// 以下在玩家登陆时，连接自己的数据库表
				db, err := ModuleManager.GetModuleManager().GetDb(Cfg.InfoDb)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					tabInfo, err := DbModule.GetInfoDb(db, DbModule.UserId(userInfo.UserId))
					if err != nil {
						fmt.Println(err.Error())
					} else {
						h.infoDb = tabInfo
					}
					// 测试
					//info := h.infoDb.GetCurStatisticalInfo()
					//fmt.Println(*info)
					//h.infoDb.Insert(DbModule.UseType("19年1月的测试"), 22)
					//infoList,err := h.infoDb.QueryStatisticalInfoList(DbModule.WEEK_INFO)
					//if err != nil {
					//	fmt.Println(err)
					//} else {
					//	fmt.Println("数据长度：", len(infoList))
					//	for _, info := range infoList{
					//		fmt.Println("信息：", info.Id, info.Time, info.UserType, info.Val)
					//	}
					//	res, _ := h.infoDb.QueryStatisticalInfo(DbModule.MON_INFO)
					//	fmt.Println(res)
					//}
				}
				h.loginRep(msg.CodeType_SUC, "登录成功", &userInfo)
			}
		}
	}
}

/**
 回复登录
 */
func (h *handle) loginRep (codeType msg.CodeType, str string, userInfo *DbModule.UserInfo)  {
	var user *msg.User
	var sInfo *msg.StatisticalInfo
	if codeType == msg.CodeType_SUC {
		user = &msg.User{Id: int32(changeId(userInfo.UserId)), Name: userInfo.UserName}
		info := h.infoDb.GetCurStatisticalInfo()
		sInfo = &msg.StatisticalInfo{
			WeekVal: info.Week,
			MonVal:  info.Mon,
			YearVal: info.Year,
		}
	}
	login := &jiZhangBo.LoginRep{
		Code: h.getCode(codeType, str),
		User: user,
		SInfo: sInfo,
	}
	h.send(msg.Event_EVENT_LOGIN_REP, login)
}

/**
 新增信息
 */
func (h *handle) addInfo (msgByte []byte)  {
	info := &jiZhangBo.AddInfoReq{}
	h.decodeData(msgByte, info)
	//fmt.Println("收到的信息：", info)
	err := h.infoDb.Insert(DbModule.UseType(info.Info.Usetype), info.Info.Val)
	if err != nil {
		fmt.Println(err.Error())
		code := h.getCode(msg.CodeType_ERR, "添加失败")
		h.send(msg.Event_EVENT_MSG_INFO, code)
	} else {
		code := h.getCode(msg.CodeType_SUC, "添加成功")
		h.send(msg.Event_EVENT_MSG_INFO, code)
		h.statisticalInfoUpdate()
	}
}

/**
 删除信息
 */
func (h *handle) delInfo (msgByte []byte)  {
	info := &jiZhangBo.DelInfoReq{}
	h.decodeData(msgByte, info)
	err := h.infoDb.DelInfo(info.Idlist)
	if err != nil {
		fmt.Println(err.Error())
		code := h.getCode(msg.CodeType_ERR, "删除失败")
		delRep := &jiZhangBo.DelInfoRep{}
		delRep.Code = code
		h.send(msg.Event_EVENT_DEL_INFO_REP, delRep)
	} else {
		code := h.getCode(msg.CodeType_SUC, "删除成功")
		delRep := &jiZhangBo.DelInfoRep{
			Code: code,
		}
		for _, id := range info.Idlist {
			delRep.Idlist = append(delRep.Idlist, id)
		}
		h.send(msg.Event_EVENT_DEL_INFO_REP, delRep)
		h.statisticalInfoUpdate()
	}
}

/**
 查询信息
 */
func (h *handle) queryInfo (msgByte []byte)  {
	queryReq := &jiZhangBo.QueryInfoReq{}
	h.decodeData(msgByte, queryReq)
	//fmt.Println(queryReq)
	infoList, err := h.infoDb.QueryInfoToTerm(DbModule.UseType(queryReq.Usetype),
		DbModule.UseTime(queryReq.StartTime), DbModule.UseTime(queryReq.EndTime))
	if err != nil {
		fmt.Println(err.Error())
		code := h.getCode(msg.CodeType_ERR, "查找失败，请重新查找")
		h.send(msg.Event_EVENT_MSG_INFO, code)
	} else {
		reInfoList := &jiZhangBo.InfoList{}
		for _, info := range infoList{
			newInfo := &msg.Info{
				Id: int32(info.Id),
				Usetype: info.UserType,
				Time: int32(info.Time),
				Val: info.Val,
			}
			reInfoList.Infolist = append(reInfoList.Infolist, newInfo)
		}
		h.send(msg.Event_EVENT_RETURN_INFO_LIST, reInfoList)
	}
}

/**
 修改信息
 */
func (h *handle) changeInfo (msgByte []byte)  {
	changeReq := &jiZhangBo.ChangeInfoReq{}
	h.decodeData(msgByte, changeReq)
	err := h.infoDb.UpdateInfo(int(changeReq.Infoid),
		DbModule.UseType(changeReq.Changeusetype), changeReq.Changeval)
	changeRep := &jiZhangBo.ChangeInfoRep{}
	code := &msg.Code{}
	if err != nil {
		fmt.Println("修改数据库信息：",err.Error())
		code.Code = msg.CodeType_ERR
		code.Msg = "修改失败"
	} else {
		info, err := h.infoDb.QueryInfoById(changeReq.Infoid)
		if err != nil {
			fmt.Println("修改后查找数据库信息：", err.Error())
			code.Code = msg.CodeType_ERR
			code.Msg = "修改成功但获取失败"
		} else {
			code.Code = msg.CodeType_SUC
			code.Msg = "修改成功"
			changeRep.Info = &msg.Info{
				Id: int32(info.Id),
				Usetype: info.UserType,
				Time: int32(info.Time),
				Val: info.Val,
			}
		}
	}
	changeRep.Code = code
	h.send(msg.Event_EVENT_CHANGE_INFO_REP, changeRep)
	h.statisticalInfoUpdate()
}

/**
 更新统计信息
 */
func (h *handle) statisticalInfoUpdate ()  {
	info := h.infoDb.GetCurStatisticalInfo()
	sInfo := &msg.StatisticalInfo{
		WeekVal: info.Week,
		MonVal:  info.Mon,
		YearVal: info.Year,
	}
	h.send(msg.Event_EVENT_STATISYICAL_INFO_CHANGE, sInfo)
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