package router

import (
	"github.com/gin-gonic/gin"
	"houze_ops_backend/api/gen_blog/controller"
)

func InitRouter(app *gin.Engine) {
	r := app.Group("/api/gen-blog")
	{
		r.GET("/category", controller.GetAllCategory())
		r.GET("", controller.GetAllBlogs())
		r.GET("/detail", controller.GetDetailBlog())
		r.POST("", controller.CreateBlog())
		r.POST("/update", controller.UpdateBlog())
		r.DELETE("", controller.DeleteBlog())
	}
}
