package controllers

import (
	"time"

	"github.com/golang-jwt/jwt"
	"wozaizhao.com/wzzapi/config"
)

var jwtSecret = config.GetConfig().JwtSecret

// Claims
type Claims struct {
	UserID uint `json:"UserID"`
	jwt.StandardClaims
}

// GenerateToken 生成token
func generateToken(UserID uint) (string, error) {
	nowtime := time.Now()
	expireTime := nowtime.Add(24 * time.Hour)
	claims := Claims{
		UserID,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "wozaizhao",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(jwtSecret))
	return token, err
}

// ParseToken 解析token
func parseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, err
		}
	}
	return nil, err
}
