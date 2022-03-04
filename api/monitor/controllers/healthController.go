package controllers

import "github.com/gin-gonic/gin"

// @BasePath /api/v1

// @Summary returns the liveness of a microservice.
// @Schemes
// @Description If the check does not return the expected response, it means that the process is unhealthy or dead and should be replaced as soon as possible.
// @Tags Monitor
// @Accept json
// @Produce json
// @Success 200 {string} Pong
// @Router /health/ping [get]
func Pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
