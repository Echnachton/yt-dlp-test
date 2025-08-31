package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Echnachton/yt-dlp-test/internal/db"
	"github.com/Echnachton/yt-dlp-test/internal/handlers"
	"github.com/Echnachton/yt-dlp-test/internal/jobManager"
	"github.com/Echnachton/yt-dlp-test/internal/logger"
	"github.com/gin-gonic/gin"
)

const WORKER_COUNT = 5

func init() {
	db.Init()
	logger.Init()
}

func main() {
	router := gin.Default()
	
	jbmgr := jobManager.NewJobManager(WORKER_COUNT)
	jbmgr.Init()
	
	handlers.SetupRoutes(router, jbmgr)
	
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