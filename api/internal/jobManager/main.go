package jobManager

import (
	"os/exec"
	"sync"

	"github.com/Echnachton/yt-dlp-test/internal/logger"
)

type Job struct {
	URL string `json:"url" binding:"required"`
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
		logger.Printf("Processing job: %s\n", job.URL)
		jobManager.processJob(job)
	}
}

func (jobManager *JobManager) processJob(job *Job) {
	cmd := exec.Command("yt-dlp", job.URL, "-P", "temp:/tmp", "-P", "home:../../videos")
	if err := cmd.Run(); err != nil {
		logger.Printf("Error downloading %s: %v\n", job.URL, err)
	} else {
		logger.Printf("Successfully downloaded %s\n", job.URL)
	}
}

func NewJobManager(workerCount int) *JobManager {
	return &JobManager{
		jobQueue:    make(chan *Job),
		waitGroup:   sync.WaitGroup{},
	}
}

func (jm *JobManager) AddJob(job *Job) {
	jm.jobQueue <- job
}

func (jm *JobManager) CloseQueue() {
	close(jm.jobQueue)
}

func (jm *JobManager) WaitForWorkers() {
	jm.waitGroup.Wait()
}