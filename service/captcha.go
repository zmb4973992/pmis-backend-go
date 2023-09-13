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
	Id           string `json:"id"`
	Base64String string `json:"base64_string"`
}

func (c *CaptchaGet) Get() (output *CaptchaOutput, errCode int) {
	height := global.Config.Captcha.ImageHeight
	width := global.Config.Captcha.ImageWidth
	length := global.Config.Captcha.DigitLength
	maxSkew := global.Config.Captcha.MaxSkew
	dotCount := global.Config.Captcha.DotCount

	store := base64Captcha.DefaultMemStore
	driver := base64Captcha.NewDriverDigit(height, width, length, maxSkew, dotCount)
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, base64String, err := captcha.Generate()
	if err != nil {
		response.GenerateCommon(nil, util.ErrorFailToGenerateCaptcha)
	}

	output.Id = id
	output.Base64String = base64String

	return output, util.Success
}
