package router

import (
	"github.com/gin-gonic/gin"
	"houze_ops_backend/api/user/controllers"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/v1/user")
	{
		r.POST("/send-otp-code", controllers.Create())
		r.POST("/resend-otp-code", controllers.Login())
	}
}
