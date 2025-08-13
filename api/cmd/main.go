package main

import (
	"sync"

	"github.com/gin-gonic/gin"
)

const WORKER_COUNT = 10;

type Job struct {
	Url string `json:"url"`
}

type JobManager struct {
	jobs []Job
	waitGroup sync.WaitGroup
	workerCount int
}

func (jobManager *JobManager) Init() {
	for i := 0; i < jobManager.workerCount; i++ {
		jobManager.waitGroup.Add(1)
		go jobManager.worker(i)
	}
}



func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	router.Run(":8080")
}