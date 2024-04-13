package middleware

import (
	"Vchat/internal/web"
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddleWear struct {
}

func NewLoginJWTMiddleWear() *LoginJWTMiddleWear {
	return &LoginJWTMiddleWear{}
}

func (w *LoginJWTMiddleWear) CheckLogin() gin.HandlerFunc {
	// 让time.Now()转为可传入Redis的字节串
	gob.Register(time.Now())
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" || path == "/users/login_sms/code/send" || path == "/users/login_sms" {
			return
		}
		authString := c.GetHeader("Authorization")
		if authString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		auths := strings.Split(authString, " ")
		if len(auths) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := auths[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			// token 解析出来了，但是 token 可能是非法的，或者过期了的
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//if uc.UserAgent != c.GetHeader("User-Agent") {
		//	// 后期我们讲到了监控告警的时候，这个地方要埋点
		//	// 能够进来这个分支的，大概率是攻击者
		//	c.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		expireTime := uc.ExpiresAt
		if expireTime.Sub(time.Now()) < time.Minute {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 5))
			tokenStr, err = token.SignedString(web.JWTKey)
			c.Header("x-jwt-token", tokenStr)
			if err != nil {
				// 这边不要中断，因为仅仅是过期时间没有刷新，但是用户是登录了的
				log.Println(err)
			}
		}
		c.Set("user", uc)
	}
}
