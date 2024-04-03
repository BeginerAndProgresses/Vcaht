package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddleWear struct {
}

func (w *LoginMiddleWear) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		if sessions.Default(c).Get("uid") == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
