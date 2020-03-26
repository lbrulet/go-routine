package jobs

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestNewNotificationJob(t *testing.T) {
	tests := []struct {
		name string
		want *NotificationJob
	}{
		{
			name: "new_notification_job",
			want: &NotificationJob{name:"notification_job", sleep: time.Duration(rand.Int()%100+10) * time.Millisecond},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewNotificationJob()
			got.Handler()
			if reflect.TypeOf(got.sleep) != reflect.TypeOf(tt.want.GetSleep()) {
				t.Errorf("NewEmailJob() = %v, want %v", got, tt.want)
			}
		})
	}
}