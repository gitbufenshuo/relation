package handler

import (
	"net/http"

	"github.com/gitbufenshuo/relation/content"

	"strconv"

	"github.com/labstack/echo"
)

// /relation/add/:self/:up
func AddHandler(c echo.Context) error {
	var self uint64
	var up uint64
	{
		ss := c.Param("self")
		if n, err := strconv.ParseUint(ss, 10, 64); err != nil {
			return c.JSON(http.StatusOK, "bad")
		} else {
			self = n
		}
	}
	{
		ss := c.Param("up")
		if n, err := strconv.ParseUint(ss, 10, 64); err != nil {
			return c.JSON(http.StatusOK, "bad")
		} else {
			up = n
		}
	}
	content.Content_lock()
	defer content.Content_unlock()
	if ok := content.Content_add(self, up); !ok {
		c.Logger().Printf("adderr_%v_%v", self, up)
		return c.JSON(http.StatusOK, "bad")
	}
	return c.JSON(http.StatusOK, "ok")
}
