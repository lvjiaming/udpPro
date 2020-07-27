package ConnectModule

import (
	"../../Cfg"
	"../../ModuleManager"
	"../DbModule"
	"fmt"
	"net"
	"os"
	"strconv"
)

var (
	UserDb *DbModule.UserDb
	userList []*handle
)

/**
 连接函数
 */

func connect(conn *net.UDPConn)  {
	//fmt.Print("有玩家连接呢\n")
	data := make([]byte, Cfg.SERVER_RECV_LEN)
	_, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		return
	}
	fmt.Println("Receive from client", addr.String(), string(data))
}

func StartServer() (err error) {
	// 这里获取到userDb
	db, err := ModuleManager.GetModuleManager().GetDb(Cfg.UserDb)
	if err == nil {
		uDb, ok := DbModule.GetUserDb(db)
		if ok != nil {
			fmt.Println(ok.Error())
		} else {
			UserDb = uDb
		}
	}
	fmt.Print("服务器开启中\n")
	adree, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(Cfg.SERVER_PORT))
	checkError(err)
	fmt.Println("监听：", adree)
	conn, err := net.ListenUDP("udp", adree)
	checkError(err)
	defer conn.Close()
	for  {
		connect(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}