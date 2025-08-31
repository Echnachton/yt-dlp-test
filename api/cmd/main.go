package main

import (
	"context"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/Echnachton/yt-dlp-test/internal/db"
	"github.com/Echnachton/yt-dlp-test/internal/jobManager"
	"github.com/Echnachton/yt-dlp-test/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const WORKER_COUNT = 5
const YOUTUBE_URL = "https://youtube.com"
var youtubeUrlRegex = regexp.MustCompile(`^https:\/\/www\.youtube\.com.*`)

func init() {
	db.Init()
	logger.Init()
}

type DownloadRequest struct {
	URL string `json:"url" binding:"required"`
}

func main() {
	router := gin.Default()
	
	apiV1Routes := router.Group("/api/v1")
	
	jbmgr := jobManager.NewJobManager(WORKER_COUNT)
	jbmgr.Init()
	
	
	apiV1Routes.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
		});
		
	apiV1Routes.POST("/download", func(c *gin.Context) {
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

		if youtubeUrlRegex.MatchString(trimmedURL) == false {
			c.JSON(400, gin.H{"message": "Invalid URL - must be a valid YouTube URL"})
			return
		}
		
		var job jobManager.Job
		job.URL = trimmedURL
		job.ID = uuid.New().String()
		job.OwnerID = c.GetHeader("test_user")
		if job.OwnerID == "" {
			job.OwnerID = "default_user" // Provide a default if header is missing
		}
		
		jbmgr.AddJob(&job)
		
		c.JSON(200, gin.H{"message": "Job queued successfully", "job_id": job.ID})
	})
	
	apiV1Routes.GET("/download", func(c *gin.Context) {
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
	})
	
	router.StaticFS("/web", gin.Dir("../web/static", true))
	
	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Println("Received shutdown signal, closing job queue...")
		jbmgr.CloseQueue()
		jbmgr.WaitForWorkers()
		logger.Println("All workers finished, shutting down...")
		cancel()
	}()

	logger.Println("Starting server on :8080...")
	
	// Run server in a goroutine so we can handle shutdown
	go func() {
		if err := router.Run(":8080"); err != nil {
			logger.Printf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Println("Server shutdown complete")
}