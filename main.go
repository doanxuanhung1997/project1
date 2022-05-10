package main

import (
	"github.com/gin-gonic/gin"
	routerUser "houze_ops_backend/api/user/router"
	"houze_ops_backend/config"
	"houze_ops_backend/db"
)

func main() {
	// environment
	arg := "dev"
	config.Loads("./" + arg + ".env")
	config.SetEnv(config.EnvData)
	env := config.GetEnvValue()
	app := gin.Default()
	routerUser.InitRouter(app)
	if err := db.InitDb(); err != nil {
		panic(err)
	}
	_ = app.Run(env.Server.Host + ":" + env.Server.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

