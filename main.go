package main

import (
	"github.com/gin-gonic/gin"
	routerUser "houze_ops_backend/api/sys_user/router"
	"houze_ops_backend/config"
	"houze_ops_backend/db"
)

func main() {
	// environment
	arg := "dev"
	config.Loads("./" + arg + ".env")
	connectDB := db.InitConnectionDB()
	defer connectDB.Close()
	config.SetEnv(config.EnvData)
	env := config.GetEnvValue()
	app := gin.Default()
	_ = app.SetTrustedProxies(nil)
	routerUser.InitRouter(app)
	_ = app.Run(env.Server.Host + ":" + env.Server.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
