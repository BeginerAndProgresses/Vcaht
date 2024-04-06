package main

import (
	"Vchat/internal/config"
	"Vchat/internal/repository"
	"Vchat/internal/repository/dao"
	"Vchat/internal/service"
	"Vchat/internal/web"
	"Vchat/internal/web/middleware"
	"Vchat/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/sessions"
	redissess "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	db = db.Debug()

	server := initWebServer()
	initUserHdl(db, server)
	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "你好")
	//})
	server.Run(":8081")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	// 跨域配置
	server.Use(middleware.NewCORS().ServeHTTP())
	//useSession(server)
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	// 限流
	server.Use(ratelimit.NewBuilder(redisClient,
		time.Second, 1).Build())
	useJWT(server)
	return server
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	uh := web.NewUserHandler(us)
	uh.RegisterRouter(server)
}

func useSession(server *gin.Engine) {
	// 存储到cookie
	//store := cookie.NewStore([]byte("secret"))
	// 存储到磁盘
	//store := memstore.NewStore([]byte("wyVoM8T8syyicWJ92kRnHb21tLqQggx1"), []byte("L0RhIUrU0uzQn5vXjnQKxX0Mf47MEzA4"))
	// 存储到redis
	store, err := redissess.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("L0RhIUrU0uzQn5vXjnQKxX0Mf47MEzB4"),
		[]byte("L0RhIUrU0uzQn5vXjnQKxX0Mf41MEzA4"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("mysession", store), new(middleware.LoginMiddleWear).CheckLogin())
}

func useJWT(server *gin.Engine) {
	server.Use(new(middleware.LoginJWTMiddleWear).CheckLogin())
}
