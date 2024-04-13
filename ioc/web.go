package ioc

import (
	"Vchat/internal/web"
	"Vchat/internal/web/middleware"
	"Vchat/pkg/ginx/middleware/ratelimit"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRouter(server)
	return server
}

func InitWebMiddleware(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.NewCORS().CORSService(),
		ratelimit.NewBuilder(redisClient, time.Second, 1000).Build(),
		middleware.NewLoginJWTMiddleWear().CheckLogin(),
	}
}
