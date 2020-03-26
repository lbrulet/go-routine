package main

import (
	"fmt"
	"github.com/lbrulet/go-routine/pkg/jobs"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getType(t *testing.T) {
	tests := []struct {
		name string
		want jobs.IJob
	}{
		{
			name: "notification_job",
			want: jobs.NewNotificationJob(),
		},
		{
			name: "email_job",
			want: jobs.NewEmailJob(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getType(); got.GetName() != tests[0].name && got.GetName() != tests[1].name {
				t.Errorf("getType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_startJobProcessor(t *testing.T) {
	type args struct {
		jobs <-chan jobs.IJob
	}
	jobsChannel := make(chan jobs.IJob, nbJobs)
	tests := []struct {
		name string
		args args
	}{
		{
			name: "testing",
			args: args{jobs: jobsChannel},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wait := startJobProcessor(tt.args.jobs)
			for i := 0; i < workers; i++ {
				wait.Done()
			}
			wait.Wait()
		})
	}
}

func Test_run(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "run test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run()
		})
	}
}

func Test_makeJob(t *testing.T) {
	tests := []struct {
		name string
		want jobs.IJob
	}{
		{
			name: "make_job_test",
			want: jobs.NewNotificationJob(),
		},
		{
			name: "make_job_test",
			want: jobs.NewEmailJob(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeJob(); got.GetName() != tests[0].want.GetName() && got.GetName() != tests[1].want.GetName() {
				t.Errorf("makeJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createJobs(t *testing.T) {
	type args struct {
		jobs chan<- jobs.IJob
	}
	jobsChannel := make(chan jobs.IJob, workers)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "create_job_test",
			args: args{jobs: jobsChannel},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createJobs(jobsChannel)
		})
	}
}

func Test_startWorker(t *testing.T) {
	type args struct {
		workerID int
		jobs     <-chan jobs.IJob
	}
	jobsChannel := make(chan jobs.IJob, workers)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "start_worker_test",
			args: args{
				workerID: 1,
				jobs:     jobsChannel,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go startWorker(tt.args.workerID, tt.args.jobs)
			job := makeJob()
			jobsChannel <- job
		})
	}
}

func TestNewServer(t *testing.T) {
	type args struct {
		jobsChannel chan jobs.IJob
	}
	jobsChannel := make(chan jobs.IJob, workers)
	mux := http.NewServeMux()
	mux.HandleFunc("/new-job", func(w http.ResponseWriter, r *http.Request) {
		job := makeJob()
		jobsChannel <- job
		_, _ = fmt.Fprintf(w, "Job %s pushed to the worker pool\n", job.GetName())
	})
	tests := []struct {
		name string
		args args
		want *http.Server
	}{
		{
			name: "new_server_test",
			args:args{jobsChannel:jobsChannel},
			want:&http.Server{
				Addr:    ":9009",
				Handler: mux,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServer(tt.args.jobsChannel)
			go got.ListenAndServe()
			_, err := http.NewRequest("GET", "localhost:9009", nil)
			if err != nil {
				t.Fatal(err)
			}
			if got.Addr != got.Addr {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler(t *testing.T) {
	jobsChannel = make(chan jobs.IJob, workers)
	req, err := http.NewRequest("GET", "localhost:9009/new-job", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	handler(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Response code was %v; want 200", res.Code)
	}
}