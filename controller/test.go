package controller

import "github.com/gin-gonic/gin"

func Test(c *gin.Context) {
	c.FileAttachment("./static/1.rar", "aaa")
}
