package router

import (
	"github.com/gitbufenshuo/relation/handler"
	"github.com/labstack/echo"
)

func Init(e *echo.Echo) {
	e.GET("/relation/add/:self/:up", handler.AddHandler)
	e.GET("/relation/up/:self", handler.AddHandler)
	e.GET("/relation/read/:what/:self", handler.ReadHandler)
}
