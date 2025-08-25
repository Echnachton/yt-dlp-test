package main

import (
	"context"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/Echnachton/yt-dlp-test/internal/jobManager"
	"github.com/Echnachton/yt-dlp-test/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const WORKER_COUNT = 5
const YOUTUBE_URL = "https://youtube.com"
var youtubeUrlRegex = regexp.MustCompile(`^https:\/\/www\.youtube\.com.*`)

func init() {
	logger.Init()
}

type DownloadRequest struct {
	URL string `json:"url" binding:"required"`
}

func main() {
	router := gin.Default()
	jbmgr := jobManager.NewJobManager(WORKER_COUNT)
	jbmgr.Init()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	});

	router.POST("/download", func(c *gin.Context) {
		var request DownloadRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			return
		}

		if youtubeUrlRegex.MatchString(request.URL) == false {
			c.JSON(400, gin.H{"message": "Invalid URL"})
			return
		}

		var job jobManager.Job
		job.URL = request.URL
		job.ID = uuid.New().String()
		job.OwnerID = c.GetHeader("test_user")
		jbmgr.AddJob(&job)

		c.JSON(200, gin.H{"message": "Job queued successfully"})
	})

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