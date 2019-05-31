package router

import (
	"github.com/gitbufenshuo/relation/handler"
	"github.com/labstack/echo"
)

func Init(e *echo.Echo) {
	e.GET("/relation/add/:self/:up", handler.AddHandler)
	e.GET("/relation/levelup/:self/:nowdeng", handler.LevelUpHandler)
	e.GET("/relation/read/:what/:self", handler.ReadHandler)
	///////////////////////////////////////////////////////
	e.GET("/relation/flushall", handler.FlushAllHandler)
	e.GET("/relation/sta", handler.StaHandler)
}
