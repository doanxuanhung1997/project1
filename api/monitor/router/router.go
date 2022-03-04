package router

import (
	"sandexcare_backend/api/monitor/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(app *gin.Engine) {
	h := app.Group("/api/v1/health")
	{
		h.GET("/ping", controllers.Pong)
		h.GET("/docs", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
