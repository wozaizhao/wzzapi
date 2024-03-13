package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"wozaizhao.com/wzzapi/common"
)

type AuthHeader struct {
	Token string `header:"Authorization"`
}

func TokenValidator(c *gin.Context) (*Claims, error) {
	h := AuthHeader{}

	if err := c.ShouldBindHeader(&h); err != nil {
		return nil, err
	}
	var token string
	if len(strings.Split(h.Token, " ")) == 2 {
		token = strings.Split(h.Token, " ")[1]
	} else {
		return nil, errors.New("token_not_found")
	}
	claims, errorParseToken := parseToken(token)
	if errorParseToken != nil {
		return nil, errorParseToken
	}
	return claims, nil
}

func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, errorTokenValidator := TokenValidator(c)
		if errorTokenValidator != nil {
			common.LogError("TokenValidator", errorTokenValidator)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", claims.UserID)

		// todo 查询用户状态,判断帐号是否已禁用
		c.Next()
	}
}
