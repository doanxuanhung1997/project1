package router

import (
	"github.com/gin-gonic/gin"
	"sandexcare_backend/api/admin/controllers"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/v1/admin")
	{
		r.GET("/withdrawal-history", controllers.GetAllWithdrawalHistory())
		r.POST("/confirm-withdrawal", controllers.ConfirmWithdrawalRequest())
		r.GET("/users", controllers.GetAllUser())
		r.POST("/users/submit-diamond", controllers.SubmitDiamondsUser())
	}
}
