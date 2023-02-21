package controller

import (
	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	c.FileAttachment("D:/1.zip", "aaa.zip")

	return
}

func Dl(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid != "" {
		c.FileAttachment("d:/1.zip", "123.zip")
		return
	}
	return
}
