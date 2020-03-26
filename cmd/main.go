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
	jobsChannel chan jobs.IJob
	workers = 0
	nbJobs = 20
)

func init() {
	flag.IntVar(&workers, "workers", 10, "Number of workers to use")
	flag.IntVar(&nbJobs, "jobs", 10, "Number of jobs to create")
}

func getType() jobs.IJob {
	return jobTypes[rand.Int()%len(jobTypes)]
}

func NewServer(jobsChannel chan jobs.IJob) *http.Server {
	addr := fmt.Sprintf(":%s", "9009")

	mux := http.NewServeMux()
	mux.HandleFunc("/new-job", handler)

	log.Println("[INFO] starting HTTP server on port :9009")
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	job := makeJob()
	jobsChannel <- job
	_, _ = fmt.Fprintf(w, "Job %s pushed to the worker pool\n", job.GetName())
}

func run() chan jobs.IJob {
	// parse the flags
	flag.Parse()

	// create a channel
	jobsChannel = make(chan jobs.IJob, nbJobs)

	// start the job processor
	go func() {
		wait := startJobProcessor(jobsChannel)
		wait.Wait()
	}()

	go createJobs(jobsChannel)
	return jobsChannel
}

// main entry point for the application
func main() {
	srv := NewServer(run())
	log.Fatal(srv.ListenAndServe())
}

// makeJob creates a new job
func makeJob() jobs.IJob {
	return getType()
}
func startJobProcessor(jobs <-chan jobs.IJob) sync.WaitGroup {
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

	return wait
}

// create an certain amount of jobs according to nbJobs and send it to the job channel
func createJobs(jobs chan<- jobs.IJob) {
	log.Printf("[INFO] starting %d jobs\n", nbJobs)
	for i := 0; i < nbJobs; i++ {
		job := makeJob()
		jobs <- job
		time.Sleep(5 * time.Millisecond)
	}
}

// creates a worker that pulls jobs from the job channel
func startWorker(workerID int, jobs <-chan jobs.IJob) {
	for {
		select {
		// read from the job channel
		case job := <-jobs:
			startTime := time.Now()
			// fake processing the request
			time.Sleep(job.GetSleep())
			log.Printf("[workerID:%d][%s] Processed job in %0.3f seconds", workerID, job.GetName(), time.Now().Sub(startTime).Seconds())
		}
	}
}
