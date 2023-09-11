package service

import (
	"github.com/mojocn/base64Captcha"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type CaptchaGet struct {
}

type CaptchaOutput struct {
	ID           string `json:"id"`
	Base64String string `json:"base64_string"`
}

func (c *CaptchaGet) Get() (output *CaptchaOutput, errCode int) {
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
		response.GenerateCommon(nil, util.ErrorFailToGenerateCaptcha)
	}

	output.ID = id
	output.Base64String = base64String

	return output, util.Success
}
