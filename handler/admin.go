package handler

import (
	"net/http"

	"github.com/gitbufenshuo/relation/content"
	"github.com/labstack/echo"
)

// 清库
// /relation/flushall
func FlushAllHandler(c echo.Context) error {
	content.Content_lock()
	defer content.Content_unlock()
	content.Content_flushall()
	return nil
}

// 查看统计 总键数，最大深度，最大宽度
// /relation/sta
func StaHandler(c echo.Context) error {
	content.Content_lock()
	defer content.Content_unlock()
	d1, d2, d3 := content.Content_sta()
	return c.JSON(http.StatusOK, ReadRes{
		Data: []uint64{d1, d2, d3},
	})
}
