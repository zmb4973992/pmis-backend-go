package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"pmis-backend-go/util/jwt"
)

type loginService struct{}

func (loginService) Login(param dto.Login) response.Common {
	var user model.User
	//根据入参的用户名，从数据库取出记录赋值给user
	err := global.DB.Where("username=?", param.Username).First(&user).Error
	//如果没有找到记录
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorInvalidUsernameOrPassword)
	}
	//如果找到记录了，但是密码错误的话
	if util.CheckPassword(param.Password, user.Password) == false {
		return response.Failure(util.ErrorInvalidUsernameOrPassword)
	}
	//账号密码都正确时，生成token
	token := jwt.GenerateToken(user.ID)
	return response.SuccessWithData(
		map[string]any{
			"access_token": token,
		})
}
