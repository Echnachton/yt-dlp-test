package handlers

import (
	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "OK"})
}

func DynamicHealthFileHandler(c *gin.Context) {
	id := c.Param("id")

	c.JSON(200, gin.H{"message": "OK", "id": id})
}