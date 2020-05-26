/**
 启动程序
 */
package main

import (
	"./Cfg"
	. "./ModuleManager"
	"./ModuleManager/ConnectModule"
	"fmt"
	"sync"
)
var waitGroup sync.WaitGroup
func main()  {
	fmt.Print("等待准备工作\n")
	waitGroup.Add(2)
	go GetModuleManager().ConnectDb(&waitGroup, Cfg.UserDb) // 在其他包调用时需要传内存地址，传值无法生效
	go GetModuleManager().ConnectDb(&waitGroup, Cfg.InfoDb) // 在其他包调用时需要传内存地址，传值无法生效
	waitGroup.Wait()
	fmt.Print("所有准备已就绪\n")
	err := ConnectModule.StartServer()
	if err != nil {
		fmt.Println(err.Error())
	}
}
