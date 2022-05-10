package controller

import (
	"github.com/gin-gonic/gin"
	"houze_ops_backend/api/sys_master/repository"
	"net/http"
)

func GetConsultingField() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := repository.PublishInterfaceMaster().GetAllConsultingField()
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": data,
		})
	}
}
