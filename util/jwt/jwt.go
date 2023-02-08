package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"pmis-backend-go/global"
	"time"
)

type CustomClaims struct {
	UserID int
	jwt.RegisteredClaims
}

var validityDays = time.Duration(global.Config.ValidityPeriod)

// 构建载荷
func buildClaims(userID int) CustomClaims {
	return CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "zhoumengbin",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(validityDays * 24 * time.Hour)),
		},
	}
}

// GenerateToken 传入userID，返回token字符串
func GenerateToken(userID int) (string, error) {
	claims := buildClaims(userID)
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenStruct.SignedString([]byte(global.Config.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 验证用户token。这部分基本就是参照官方写法。
// 第一个参数是token字符串，第二个参数是结构体，第三个参数是jwt规定的解析函数，包含密钥
func ParseToken(token string) (*CustomClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return global.Config.SecretKey, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*CustomClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
