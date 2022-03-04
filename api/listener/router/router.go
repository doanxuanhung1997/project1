package router

import (
	"sandexcare_backend/api/listener/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/v1/listener")
	{
		r.POST("/login", controllers.Login())
		r.POST("/login2", controllers.Login2())
		r.GET("/verify-token", controllers.VerifyToken())
		
		r.POST("/logout", controllers.Logout())
		r.POST("/create", controllers.CreateListener())
		r.POST("/forgot-password", controllers.ForgotPassword())
		r.POST("/resend-forgot-password", controllers.ResendForgotPassword())
		r.POST("/verify-reset-password", controllers.VerifyResetPassword())
		r.POST("/reset-new-password", controllers.ResetNewPassword())
		r.POST("/change-password", controllers.ChangePassword())
		r.GET("/info", controllers.GetInfoListener())
		r.POST("/request-withdrawal", controllers.RequestWithdrawal())
		r.GET("/withdrawal-history", controllers.GetWithdrawalHistory())
		r.GET("/revenue-analysis", controllers.GetRevenueAnalysis())
		r.POST("/miss-call", controllers.HandleMissCall())
	}
}
