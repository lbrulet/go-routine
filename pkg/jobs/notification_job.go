package jobs

import (
	"fmt"
	"math/rand"
	"time"
)

type NotificationJob struct {
	name string
	sleep time.Duration
}

func NewNotificationJob() *NotificationJob {
	return &NotificationJob{
		name: "notification_job",
		sleep: time.Duration(rand.Int()%100+10) * time.Millisecond,
	}
}

func (n *NotificationJob) GetName() string {
	return n.name
}

func (n *NotificationJob) GetSleep() time.Duration {
	return n.sleep
}

func (n *NotificationJob) Handler() {
	fmt.Println(n.name)
}
