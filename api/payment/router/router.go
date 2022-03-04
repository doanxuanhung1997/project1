package router

import (
	"github.com/gin-gonic/gin"
	"sandexcare_backend/api/payment/controllers"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/v1/payment")
	{
		r.POST("/", controllers.OrderPayment())
		r.POST("/booking-expert", controllers.PaymentBookingExpert())
		r.POST("/call", controllers.CallPayment())
		r.POST("/update", controllers.UpdateOrderPayment())
		r.POST("/refund", controllers.PaymentRefund())
		r.POST("/refund-expert", controllers.PaymentRefundExpert())
		r.POST("/unlock", controllers.Unlock())
	}
}
