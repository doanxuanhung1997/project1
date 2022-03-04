package router

import (
	"github.com/gin-gonic/gin"
	"sandexcare_backend/api/track_event/controllers"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/v1/track-event")
	{
		r.POST("/listener/action", controllers.TrackingActionListener())
	}
}
