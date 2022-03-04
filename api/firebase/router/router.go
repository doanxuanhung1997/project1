package router

import (
	"sandexcare_backend/api/firebase/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter(app *gin.Engine) {
	h := app.Group("/api/v1/")
	{
		h.POST("/firebase/send", controllers.SendNotify())
		h.POST("/firebase/register", controllers.RegisterNotify())
	}
}
