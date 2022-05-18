package router

import (
	"github.com/gin-gonic/gin"
	"houze_ops_backend/api/sys_user/controller"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("api/sys-user")
	{
		r.POST("/login", controller.Login())
		r.POST("/create", controller.CreateUser())
	}
}
