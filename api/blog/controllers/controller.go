package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"houze_ops_backend/api/master/repository"
)

func GetConsultingField() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := repository.PublishInterfaceMaster().GetAllConsultingField()
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"data":    data,
		})
	}
}