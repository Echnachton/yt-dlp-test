package handlers

import (
	"regexp"
	"strings"

	"github.com/Echnachton/yt-dlp-test/internal/db"
	"github.com/Echnachton/yt-dlp-test/internal/jobManager"
	"github.com/Echnachton/yt-dlp-test/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var youtubeUrlRegex = regexp.MustCompile(`^https:\/\/www\.youtube\.com.*`)

type DownloadRequest struct {
	URL string `json:"url" binding:"required"`
}

func DownloadHandler(jbmgr *jobManager.JobManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request DownloadRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"message": "Invalid JSON", "error": err.Error()})
			return
		}

		trimmedURL := strings.TrimSpace(request.URL)
		if trimmedURL == "" {
			c.JSON(400, gin.H{"message": "URL cannot be empty"})
			return
		}

		if !youtubeUrlRegex.MatchString(trimmedURL) {
			c.JSON(400, gin.H{"message": "Invalid URL - must be a valid YouTube URL"})
			return
		}
		
		var job jobManager.Job
		job.URL = trimmedURL
		job.ID = uuid.New().String()
		job.OwnerID = c.GetHeader("test_user")
		if job.OwnerID == "" {
			job.OwnerID = "default_user"
		}
		
		jbmgr.AddJob(&job)
		
		c.JSON(200, gin.H{"message": "Job queued successfully", "job_id": job.ID})
	}
}

func GetDownloadsHandler(c *gin.Context) {
	queryResult, err := db.GetDB().Query("SELECT id, url, owner_id FROM videos")
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}
	defer queryResult.Close()
	
	var videos []jobManager.Job
	for queryResult.Next() {
		var video jobManager.Job
		if err := queryResult.Scan(&video.ID, &video.URL, &video.OwnerID); err != nil {
			logger.Printf("Error scanning row: %v", err)
			continue
		}
		videos = append(videos, video)
	}
	c.JSON(200, gin.H{"videos": videos})
}
