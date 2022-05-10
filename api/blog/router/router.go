package router

import (
	"github.com/gin-gonic/gin"
	"houze_ops_backend/api/master/controllers"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/v1/master")
	{
		r.GET("/consulting-field", controllers.GetConsultingField())
		r.GET("/time-slot", controllers.GetTimeSlot())
		r.GET("/config", controllers.GetConfigData())
		r.POST("/send-event-ws", controllers.SendEventWebSocket())
		r.GET("/datetime", controllers.GetServerDatetime())
	}
}
