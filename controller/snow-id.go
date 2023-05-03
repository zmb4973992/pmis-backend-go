package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type snowID struct{}

func (s *snowID) Get(c *gin.Context) {
	snowID, err := util.Snowflake.NextID()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorFailToGenerateSnowID))
		return
	}

	c.JSON(http.StatusOK, snowID)
	return
}
