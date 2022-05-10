package controllers

import (
	"github.com/gin-gonic/gin"
	"houze_ops_backend/api/user/repository"
	"net/http"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.JSON(http.StatusOK, nil)
		return
	}
}

func Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := repository.PublishInterfaceUser().CreateUser();
		c.JSON(http.StatusOK, data)
		return
	}
}
