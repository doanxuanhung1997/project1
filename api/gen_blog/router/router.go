package router

import (
	"github.com/gin-gonic/gin"
	"houze_ops_backend/api/sys_master/controller"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/v1/sys_master")
	{
		r.GET("/consulting-field", controller.GetConsultingField())
		r.GET("/time-slot", controller.GetTimeSlot())
		r.GET("/config", controller.GetConfigData())
		r.POST("/send-event-ws", controller.SendEventWebSocket())
		r.GET("/datetime", controller.GetServerDatetime())
	}
}
