package router

import (
	"github.com/gin-gonic/gin"
	"sandexcare_backend/api/conversation/controllers"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/v1/conversation")
	{
		// user request API
		r.POST("/start-call", controllers.StartConversation())
		r.POST("/start-call/expert", controllers.StartCallExpert())
		r.POST("/end-call", controllers.EndCall())
		r.POST("/switch-listener", controllers.SwitchListener())
		r.POST("/evaluation", controllers.SubmitCallEvaluation())
		r.GET("/user", controllers.GetConversationHistoryUser())
		r.GET("/user/detail", controllers.GetDetailConversationHistoryUser())
		r.POST("/accept-extend-call", controllers.AcceptExtendCall())

		// listener request API
		r.POST("/join-call", controllers.JoinConversation())
		r.POST("/info-call", controllers.SubmitInfoConversation())
		r.POST("/request-extend-call", controllers.RequestExtendCall())
		r.GET("/call-history/user", controllers.GetCallHistoryUser())
		r.GET("/call-history/listener", controllers.GetCallHistoryListener())
		r.GET("/call-history/listener/detail", controllers.GetDetailCallHistoryForListener())

	}
}
