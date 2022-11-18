package jwt

import (
	"github.com/golang-jwt/jwt"
	"pmis-backend-go/global"
	"time"
)

type MyClaim struct {
	UserID int
	jwt.StandardClaims
}

// GenerateToken 传入Username，返回token字符串
func GenerateToken(userID int) string {
	days := time.Duration(global.Config.ValidityPeriod)
	claim := MyClaim{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(days * 24 * time.Hour).Unix(),
		}}
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, _ := tokenStruct.SignedString(global.Config.SecretKey)
	return tokenString
}

// ParseToken 验证用户token。这部分基本就是参照官方写法。
//第一个参数是token字符串，第二个参数是结构体，第三个参数是jwt规定的解析函数，包含密钥
func ParseToken(token string) (*MyClaim, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &MyClaim{}, func(token *jwt.Token) (interface{}, error) {
		return global.Config.SecretKey, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*MyClaim); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
