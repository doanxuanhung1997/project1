package router

import (
	"github.com/gin-gonic/gin"
	"houze_ops_backend/api/sys_master/controller"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/sys-master")
	{
		r.GET("/province", controller.GetProvince())
		r.GET("/district", controller.GetDistrict())
		r.GET("/wards", controller.GetWards())
	}
}
