package controller

import (
	"github.com/gin-gonic/gin"
	sysMasterRepo "houze_ops_backend/api/sys_master/repository"
	"houze_ops_backend/helpers/message"
	"net/http"
)

func GetProvince() gin.HandlerFunc {
	return func(c *gin.Context) {
		province := sysMasterRepo.PublishInterface().GetProvince()
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
			"data":    province,
		})
		return
	}
}

func GetDistrict() gin.HandlerFunc {
	return func(c *gin.Context) {
		parentCode := c.DefaultQuery("parent_code", "79")
		district := sysMasterRepo.PublishInterface().GetDistrict(parentCode)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
			"data":    district,
		})
		return
	}
}
func GetWards() gin.HandlerFunc {
	return func(c *gin.Context) {
		parentCode := c.DefaultQuery("parent_code", "760")
		wards := sysMasterRepo.PublishInterface().GetWards(parentCode)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
			"data":    wards,
		})
		return
	}
}
