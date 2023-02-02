package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"pmis-backend-go/util/jwt"
)

type login struct{}

func (*login) Login(param dto.Login) response.Common {
	var record model.User
	//根据入参的用户名，从数据库取出记录赋值给user
	err := global.DB.Where("username=?", param.Username).First(&record).Error
	//如果没有找到记录
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorInvalidUsernameOrPassword)
	}
	//如果密码错误
	if util.CheckPassword(param.Password, record.Password) == false {
		return response.Fail(util.ErrorInvalidUsernameOrPassword)
	}
	//账号密码都正确时，生成token
	token := jwt.GenerateToken(record.ID)
	return response.SucceedWithData(
		map[string]any{
			"access_token": token,
		})
}
