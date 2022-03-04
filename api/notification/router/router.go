package router

import (
	"sandexcare_backend/api/notification/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/v1/notification")
	{
		r.POST("/", controllers.CreateNotification())
		r.GET("/", controllers.GetListNotification())
		r.POST("/read", controllers.ReadNotification())
	}
}
