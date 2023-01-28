package util

import (
	"golang.org/x/crypto/bcrypt"
	"pmis-backend-go/global"
)

func EncryptPassword(originalPassword string) (encryptedPassword string, err error) {
	//10为加密难度，取值范围为4-31，官方建议10
	bytes, err := bcrypt.GenerateFromPassword([]byte(originalPassword), 10)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return "", err
	}
	encryptedPassword = string(bytes)
	return encryptedPassword, nil
}

func CheckPassword(originalPassword string, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(originalPassword))
	//等于如果err有值，那么err==nil为false，返回false；如果err为空，err==nil为true，返回true
	return err == nil
}
