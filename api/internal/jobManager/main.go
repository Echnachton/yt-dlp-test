package jobManager

import (
	"bytes"
	"os/exec"
	"sync"

	"github.com/Echnachton/yt-dlp-test/internal/db"
	"github.com/Echnachton/yt-dlp-test/internal/logger"
)

type Job struct {
	ID string `json:"id"`
	URL string `json:"url" binding:"required"`
	OwnerID string `json:"owner_id" binding:"required"`
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
	db.GetDB().Exec("INSERT INTO videos (url, internal_video_id, owner_id, status) VALUES (?, ?, ?, ?)", job.URL, job.ID, job.OwnerID, "PENDING")
	outDir :="home:../../videos/" + job.ID
	cmd := exec.Command(
		"yt-dlp", job.URL, 
	"-P", "temp:/tmp", 
	"-P", outDir)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logger.Printf("Error downloading %s: %s\n", job.URL, err)
		
		if stdout.Len() > 0 {
			logger.Printf("Stdout: %s\n", stdout.String())
		}

		if stderr.Len() > 0 {
			logger.Printf("  Stderr: %s\n", stderr.String())
		}
		
		db.GetDB().Exec("UPDATE videos SET status = 'FAILED' WHERE internal_video_id = ?", job.ID)
	} else {
		logger.Printf("Successfully downloaded %s (Job ID: %s)\n", job.URL, job.ID)
		db.GetDB().Exec("UPDATE videos SET status = 'COMPLETED' WHERE internal_video_id = ?", job.ID)
	}
}

func NewJobManager(workerCount int) *JobManager {
	return &JobManager{
		jobQueue:    make(chan *Job, 100),
		waitGroup:   sync.WaitGroup{},
		workerCount: workerCount,
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