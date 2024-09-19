package main

import (
	"testing"
	"time"
)

func TestScribeMetrics(t *testing.T) {

	m := metrics{}

	tests := []struct {
		name          string
		met           *metrics
		pollInterval  time.Duration
		stop          int64
		wantpoolcount int64
	}{
		{
			name:          "test subscribe #1",
			met:           &m,
			pollInterval:  1,
			stop:          1,
			wantpoolcount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ScribeMetrics(tt.met, tt.pollInterval, tt.stop)
			if tt.met.PollCount != tt.wantpoolcount {
				t.Errorf("ScribeMetrics() Pollcount = %d, WantPoolCount %d", tt.met.PollCount, tt.wantpoolcount)
			}
		})
	}
}
