package service

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type CaptchaGet struct {
}

func (c *CaptchaGet) Get() response.Common {
	height := global.Config.CaptchaConfig.ImageHeight
	width := global.Config.CaptchaConfig.ImageWidth
	length := global.Config.CaptchaConfig.DigitLength
	maxSkew := global.Config.CaptchaConfig.MaxSkew
	dotCount := global.Config.CaptchaConfig.DotCount

	store := base64Captcha.DefaultMemStore
	driver := base64Captcha.NewDriverDigit(height, width, length, maxSkew, dotCount)
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, base64String, err := captcha.Generate()
	if err != nil {
		global.SugaredLogger.Errorln("生成验证码失败")
		response.Failure(util.ErrorFailToGenerateCaptcha)
	}

	return response.SuccessWithData(
		gin.H{
			"id":            id,
			"base64_string": base64String,
		})
}
