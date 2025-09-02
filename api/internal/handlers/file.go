package handlers

import (
	"os"

	"github.com/Echnachton/yt-dlp-test/internal/db"
	"github.com/Echnachton/yt-dlp-test/internal/logger"
	"github.com/gin-gonic/gin"
)

func GetFileHandler(c *gin.Context) {
	var path string
	id := c.Param("id")
	queryResults, err := db.GetDB().Query("SELECT internal_video_id FROM videos WHERE internal_video_id = ?", id)
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

	c.File("../../videos/" + path + fileName)
}