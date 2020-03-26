package jobs

import (
	"fmt"
	"math/rand"
	"time"
)

type EmailJob struct {
	name string
	sleep time.Duration
}

func NewEmailJob() *EmailJob {
	return &EmailJob{
		name: "email_job",
		sleep: time.Duration(rand.Int()%100+10) * time.Millisecond,
	}
}

func (e *EmailJob) GetName() string {
	return e.name
}

func (e *EmailJob) GetSleep() time.Duration {
	return e.sleep
}

func (e *EmailJob) Handler() {
	fmt.Println(e.name)
}