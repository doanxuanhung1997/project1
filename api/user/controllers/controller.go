package controllers

import (
	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusOK, nil)
		return
	}
}