package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gitbufenshuo/relation/content"
	"github.com/gitbufenshuo/relation/router"
	"github.com/labstack/echo"
)

func main() {
	{
		content.Content_Init()
	}
	{
		go echoInit()
	}

	{ // 优雅的退出
		c := make(chan os.Signal)
		// SIGTERM 结束程序
		// 捕捉(ctrl+c)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		<-c
		// 优雅的退出
		fmt.Println("安全且优雅的退出")
	}
}

func echoInit() {
	// echo 初始化
	e := echo.New()
	router.Init(e)
	e.Start(os.Getenv("HTTP_ADDRESS") + ":" + os.Getenv("HTTP_PORT"))
}
