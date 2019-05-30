package handler

import (
	"net/http"

	"strconv"

	"github.com/gitbufenshuo/relation/content"
	"github.com/labstack/echo"
)

type ReadRes struct {
	Msg     string
	Data    []uint64
	OneData uint64
}

func ReadHandler(c echo.Context) error {
	var self uint64
	{
		ss := c.Param("self")
		if n, err := strconv.ParseUint(ss, 10, 64); err != nil {
			return c.JSON(http.StatusOK, ReadRes{
				Msg: "bad",
			})
		} else {
			self = n
		}
	}
	switch c.Param("what") {
	case "allup":
		content.Content_lock()
		defer content.Content_unlock()
		rs := content.Content_getallup(self)
		return c.JSON(http.StatusOK, ReadRes{
			Data: rs,
		})
	case "oneup":
		content.Content_lock()
		defer content.Content_unlock()
		rs := content.Content_getoneup(self)
		return c.JSON(http.StatusOK, ReadRes{
			OneData: rs,
		})
	default:
		return c.JSON(http.StatusOK, ReadRes{
			Msg: "bad",
		})
	}
}
