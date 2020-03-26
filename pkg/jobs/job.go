package jobs

import "time"

// interface to represent a job
type IJob interface {
	Handler()
	GetName() string
	GetSleep() time.Duration
}