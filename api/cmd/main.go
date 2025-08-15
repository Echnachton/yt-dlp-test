package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
)

const WORKER_COUNT = 5
const YOUTUBE_URL = "https://youtube.com"
var youtubeUrlRegex = regexp.MustCompile(`^https:\/\/www\.youtube\.com.*`)

var (
	outfile *os.File
	logger  *log.Logger
)

func init() {
	// https://specifications.freedesktop.org/basedir-spec/latest/index.html#variables
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	logDir := filepath.Join(homeDir, ".local", "state")

	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create log dir: %v", err)
	}

	outfile, err = os.Create(filepath.Join(logDir, "yt-dlp.log"))
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logger = log.New(outfile, "", log.LstdFlags|log.Lshortfile)
}

type Job struct {
	Url string `json:"url" binding:"required"`
}

type JobManager struct {
	jobQueue chan *Job
	waitGroup sync.WaitGroup
	workerCount int
}

func (jobManager *JobManager) Init() {
	for i := 0; i < jobManager.workerCount; i++ {
		jobManager.waitGroup.Add(1)
		go jobManager.worker()
	}
}

func (jobManager *JobManager) worker() {
	defer jobManager.waitGroup.Done()

	for job := range jobManager.jobQueue {
		logger.Printf("Processing job: %s\n", job.Url)
		jobManager.processJob(job)
	}
}

func (jobManager *JobManager) processJob(job *Job) {
	cmd := exec.Command("yt-dlp", job.Url, "-P", "temp:/tmp", "-P", "home:../../videos")
	if err := cmd.Run(); err != nil {
		logger.Printf("Error downloading %s: %v\n", job.Url, err)
	} else {
		logger.Printf("Successfully downloaded %s\n", job.Url)
	}
}

func NewJobManager() *JobManager {
	return &JobManager{
		jobQueue:    make(chan *Job),
		waitGroup:   sync.WaitGroup{},
		workerCount: WORKER_COUNT,
	}
}

func main() {
	// Ensure file is closed when program exits
	defer func() {
		if outfile != nil {
			logger.Println("Shutting down server...")
			outfile.Close()
		}
	}()

	router := gin.Default()
	jobManager := NewJobManager()
	jobManager.Init()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	});

	router.POST("/download", func(c *gin.Context) {
		var request Job
		if err := c.ShouldBindJSON(&request); err != nil {
			return
		}

		if youtubeUrlRegex.MatchString(request.Url) == false {
			c.JSON(400, gin.H{"message": "Invalid URL"})
			return
		}

		jobManager.jobQueue <- &request

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
		close(jobManager.jobQueue)
		jobManager.waitGroup.Wait()
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