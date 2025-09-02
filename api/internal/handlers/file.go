package handlers

import (
	"os"

	"github.com/Echnachton/yt-dlp-test/internal/db"
	"github.com/Echnachton/yt-dlp-test/internal/logger"
	"github.com/gin-gonic/gin"
)

type FileRequest struct {
	ID string `json:"id" binding:"required"`
}

func PostFileHandler(c *gin.Context) {
	var request FileRequest
	var path string

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"message": "Invalid JSON", "error": err.Error()})
		return
	}
	
	queryResults, err := db.GetDB().Query("SELECT internal_video_id FROM videos WHERE internal_video_id = ?", request.ID)
	defer queryResults.Close()

	if err != nil {
		logger.Printf("Error getting file: %v\n", err)
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	for queryResults.Next() {
		if err := queryResults.Scan(&path); err != nil {
			logger.Printf("Error scanning row: %v\n", err)
			c.JSON(404, gin.H{"message": "File not found"})
			return
		}
	}

	fileInfo, err := os.ReadDir("../../videos/" + path)
	if err != nil {
		logger.Printf("Error reading directory: %v\n", err)
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	fileName := fileInfo[0].Name()

	c.File("../../videos/" + path + "/" + fileName)
}