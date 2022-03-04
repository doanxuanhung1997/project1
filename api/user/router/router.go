package router

import (
	"sandexcare_backend/api/user/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/v1/user")
	{
		r.POST("/send-otp-code", controllers.SendOTPCode())
		r.POST("/resend-otp-code", controllers.ResendOTPCode())
		r.POST("/verify-otp-code", controllers.VerifyOTPCode())
		r.POST("/verify-otp-code2", controllers.VerifyOTPCode2())
		r.GET("/verify-token", controllers.VerifyToken())
		r.POST("/complete-info", controllers.CompleteInfo())
		r.GET("/info", controllers.GetInfoUser())
		r.POST("/update-info", controllers.UpdateUserInfo())
		r.POST("/listeners-bookmark", controllers.CreateListenersBookmark())
		r.GET("/listeners-bookmark", controllers.GetListenersBookmark())
		r.DELETE("/listeners-bookmark", controllers.DeleteListenersBookmark())
		r.GET("/coupons", controllers.GetCouponsUser())
		r.GET("/appointment-schedule", controllers.GetAppointmentScheduleForUser())
		r.GET("/appointment-schedule/detail", controllers.GetDetailAppointmentSchedule())
	}
}
