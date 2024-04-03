package main

import (
	"Vchat/internal/repository"
	"Vchat/internal/repository/dao"
	"Vchat/internal/service"
	"Vchat/internal/web"
	"Vchat/internal/web/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	db = db.Debug()

	server := initWebServer()
	initUserHdl(db, server)

	server.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/v_chat"))
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
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store), new(middleware.LoginMiddleWear).CheckLogin())
	return server
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	uh := web.NewUserHandler(us)
	uh.RegisterRouter(server)
}
