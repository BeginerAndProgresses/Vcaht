package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddleWear struct {
}

func (w *LoginMiddleWear) CheckLogin() gin.HandlerFunc {
	// 让time.Now()转为可传入Redis的字节串
	gob.Register(time.Now())
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		session := sessions.Default(c)
		uid := session.Get("uid")
		if uid == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		const updateTimeKey = "update_time"
		timeVal := session.Get(updateTimeKey)
		t, ok := timeVal.(time.Time)
		if timeVal == nil || !ok || now.Sub(t) > time.Minute*30 {
			//	更新过期时间
			session.Set(updateTimeKey, now)
			session.Set("uid", uid)
			err := session.Save()
			if err != nil {
				//		日志
				fmt.Println("-----")
			}
		}
	}
}
