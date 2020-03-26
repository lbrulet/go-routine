package main

import (
	"flag"
	"fmt"
	"github.com/lbrulet/go-routine/pkg/jobs"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

)

var (
	jobTypes = []jobs.IJob{jobs.NewEmailJob(), jobs.NewNotificationJob()}
	workers = 0
	nbJobs = 20
)

func init() {
	flag.IntVar(&workers, "workers", 10, "Number of workers to use")
	flag.IntVar(&nbJobs, "jobs", 10, "Number of jobs to create")
}

func getType() *jobs.IJob {
	return &jobTypes[rand.Int()%len(jobTypes)]
}

// main entry point for the application
func main() {
	// parse the flags
	flag.Parse()

	// create a channel
	jobsChannel := make(chan *jobs.IJob, nbJobs)

	// start the job processor
	go startJobProcessor(jobsChannel)

	go createJobs(jobsChannel)

	handler := http.NewServeMux()

	handler.HandleFunc("/new-job", func(w http.ResponseWriter, r *http.Request) {
		job := makeJob()
		jobsChannel <- job
		_, _ = fmt.Fprintf(w, "Job %s pushed to the worker pool\n", (*job).GetName())
	})

	log.Println("[INFO] starting HTTP server on port :9009")
	log.Fatal(http.ListenAndServe(":9009", handler))
}

// makeJob creates a new job
func makeJob() *jobs.IJob {
	return getType()
}

func startJobProcessor(jobs <-chan *jobs.IJob) {
	log.Printf("[INFO] starting %d workers\n", workers)
	wait := sync.WaitGroup{}
	wait.Add(workers)

	// start works according to worker variable
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			// start the worker
			startWorker(workerID, jobs)
			wait.Done()
		}(i)
	}

	wait.Wait()
}

// create an certain amount of jobs according to nbJobs and send it to the job channel
func createJobs(jobs chan<- *jobs.IJob) {
	log.Printf("[INFO] starting %d jobs\n", nbJobs)
	for i := 0; i < nbJobs; i++ {
		job := makeJob()
		jobs <- job
		time.Sleep(5 * time.Millisecond)
	}
}

// creates a worker that pulls jobs from the job channel
func startWorker(workerID int, jobs <-chan *jobs.IJob) {
	for {
		select {
		// read from the job channel
		case job := <-jobs:
			startTime := time.Now()
			// fake processing the request
			time.Sleep((*job).GetSleep())
			log.Printf("[workerID:%d][%s] Processed job in %0.3f seconds", workerID, (*job).GetName(), time.Now().Sub(startTime).Seconds())
		}
	}
}
