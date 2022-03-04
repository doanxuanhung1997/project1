package router

import (
	"github.com/gin-gonic/gin"
	"sandexcare_backend/api/schedule/controllers"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/v1/schedule-work")
	{
		r.GET("/", controllers.GetScheduleWorkForListener())
		r.GET("/day", controllers.GetScheduleInDay())
		r.POST("/create", controllers.CreateSchedule())
		r.GET("/detail-date/listener", controllers.GetDetailScheduleWorkListener())

		r.GET("/appointment", controllers.GetScheduleWorkAppointment())
		r.GET("/appointment/listener", controllers.GetListenerForBookAppointment())
		//r.GET("/appointment/expert", controllers.GetDetailScheduleWorkExpertForBookAppointment())

	}
}
