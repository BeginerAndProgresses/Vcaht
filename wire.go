//go:build wireinject

package main

import (
	"Vchat/internal/repository"
	"Vchat/internal/repository/cache"
	"Vchat/internal/repository/dao"
	"Vchat/internal/service"
	"Vchat/internal/web"
	"Vchat/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		// DAO 部分
		dao.NewUserDao,

		// cache 部分
		cache.NewRedisUserCache, cache.NewRedisCodeCache,

		// repository 部分
		repository.NewUserRepository,
		repository.NewCacheCodeRepository,

		// Service 部分
		ioc.InitSmsService,
		service.NewUserService,
		service.NewCodeService,

		// handler 部分
		web.NewUserHandler,

		ioc.InitWebMiddleware,
		ioc.InitWebServer,
	)
	return gin.Default()
}
