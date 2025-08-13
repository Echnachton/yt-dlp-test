package main

import (
	"fmt"
	"regexp"
	"sync"

	"os/exec"

	"github.com/gin-gonic/gin"
)

const WORKER_COUNT = 5
const YOUTUBE_URL = "https://youtube.com"

var youtubeUrlRegex = regexp.MustCompile(`^https:\/\/www\.youtube\.com.*`)

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
		fmt.Println("Processing job", job.Url)
		jobManager.processJob(job)
	}
}

func (jobManager *JobManager) processJob(job *Job) {
	cmd := exec.Command("yt-dlp", job.Url)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error downloading %s: %v\n", job.Url, err)
	} else {
		fmt.Printf("Successfully downloaded %s\n", job.Url)
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

	router.Run(":8080")
}