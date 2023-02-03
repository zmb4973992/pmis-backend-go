package util

import (
	"golang.org/x/crypto/bcrypt"
	"pmis-backend-go/global"
)

func Encrypt(original string) (encrypted string, err error) {
	//第二个参数为加密难度，取值范围为4-31，官方建议10。值越大，越占用cpu
	bytes, err := bcrypt.GenerateFromPassword([]byte(original), 4)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return "", err
	}
	encrypted = string(bytes)
	return encrypted, nil
}

func CheckPassword(originalPassword string, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(originalPassword))
	//等于如果err有值，那么err==nil为false，返回false；如果err为空，err==nil为true，返回true
	return err == nil
}
