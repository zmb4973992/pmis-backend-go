package util

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"pmis-backend-go/global"
	"time"
)

type CustomClaims struct {
	UserId int64
	jwt.RegisteredClaims
}

// 构建载荷
func buildClaims(userId int64) CustomClaims {
	validityDays := time.Duration(global.Config.Jwt.ValidityDays) * 24 * time.Hour
	return CustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    global.Config.Jwt.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(validityDays)),
		},
	}
}

// GenerateToken 传入userId，返回token字符串
func GenerateToken(userId int64) (string, error) {
	claims := buildClaims(userId)
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenStruct.SignedString([]byte(global.Config.Jwt.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 验证用户token。这部分基本就是参照官方写法。
// 第一个参数是token字符串，第二个参数是结构体，第三个参数是jwt规定的解析函数，包含密钥
func ParseToken(token string) (*CustomClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &CustomClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(global.Config.Jwt.SecretKey), nil
		})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*CustomClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// GetUserId 从token获取userId
func GetUserId(c *gin.Context) (userId int64, exists bool) {
	token := c.GetHeader("access_token")
	if token == "" {
		return 0, false
	}
	//开始校验access_token
	customClaims, err := ParseToken(token)
	//如果存在错误或token已过期
	if err != nil || customClaims.ExpiresAt.Unix() < time.Now().Unix() {
		return 0, false
	}
	//如果access_token校验通过
	userId = customClaims.UserId
	return userId, true
}
