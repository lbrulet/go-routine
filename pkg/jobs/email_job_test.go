package jobs

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestNewEmailJob(t *testing.T) {
	tests := []struct {
		name string
		want *EmailJob
	}{
		{
			name: "new_email_job",
			want: &EmailJob{name:"email_job", sleep: time.Duration(rand.Int()%100+10) * time.Millisecond},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEmailJob()
			got.Handler()
			if reflect.TypeOf(got.sleep) != reflect.TypeOf(tt.want.GetSleep()) {
				t.Errorf("NewEmailJob() = %v, want %v", got, tt.want)
			}
		})
	}
}