package main

import (
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	routerGenBlog "houze_ops_backend/api/gen_blog/router"
	routerSysMaster "houze_ops_backend/api/sys_master/router"
	routerUser "houze_ops_backend/api/sys_user/router"
	"houze_ops_backend/configs"
	"houze_ops_backend/db"
	"time"
)

func main() {
	// init env config
	configs.InitEnvConfig()
	env := configs.GetEnvConfig()
	// init connect database
	connectDB := db.InitConnectionDB()
	defer connectDB.Close()
	app := gin.Default()
	_ = app.SetTrustedProxies(nil)
	app.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))
	routerUser.InitRouter(app)
	routerSysMaster.InitRouter(app)
	routerGenBlog.InitRouter(app)
	_ = app.Run(env.ServerHost + ":" + env.ServerPort) // listen and serve on 0.0.0.0:2997 (for windows "localhost:2997")
}
