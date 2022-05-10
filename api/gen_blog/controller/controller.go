package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetConsultingField() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
		})
	}
}
