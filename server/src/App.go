/**
 启动程序
 */
package main

import (
	"./Cfg"
	"./Gee"
	. "./ModuleManager"
	"./ModuleManager/ConnectModule"
	"./ws"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)
var waitGroup sync.WaitGroup
func main()  {
	// Gin框架运用websocket
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to gin")
	})
	wsGroup := router.Group("/ws")
	wsGroup.GET("/:channel", ws.WebsocketManager.WsClient)
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Start Error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<- quit
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Error:", err)
	}
	log.Println("Server Shutdown")

	return
	// Gee框架
	gee := Gee.New()
	gee.GET("/", func(c *Gee.Context) {
		fmt.Println(c.Method + "-" +c.Path)
	})
	gee.GET("/hello/:name", func(c *Gee.Context) {
		fmt.Println(c.Method + "-" +c.Params["name"])
	})
	gee.GET("/ass/*filepath", func(c *Gee.Context) {
		fmt.Println(c.Method + "-" +c.Path)
	})
	runErr := gee.Run(":9999")
	if runErr != nil {
		fmt.Println(runErr.Error())
	}
	return
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
