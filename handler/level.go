package handler

import (
	"net/http"
	"strconv"

	"github.com/gitbufenshuo/relation/content"
	"github.com/labstack/echo"
)

// /relation/up/:self
func LevelUpHandler(c echo.Context) error {
	var self uint64
	{
		ss := c.Param("self")
		if n, err := strconv.ParseUint(ss, 10, 64); err != nil {
			return c.JSON(http.StatusOK, "bad")
		} else {
			self = n
		}
	}
	content.Content_lock()
	defer content.Content_unlock()
	if res := content.Content_levelup(self); res {
		return c.JSON(http.StatusOK, ReadRes{})
	} else {
		return c.JSON(http.StatusOK, ReadRes{
			Msg: "bad",
		})
	}
}
