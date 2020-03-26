package jobs

import "time"

type IJob interface {
	Handler()
	GetName() string
	GetSleep() time.Duration
}
