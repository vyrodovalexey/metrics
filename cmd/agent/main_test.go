package main

import (
	"testing"
	"time"
)

func TestScribeMetrics(t *testing.T) {

	m := metrics{}

	tests := []struct {
		name           string
		met            *metrics
		pollInterval   time.Duration
		stopcount      int64
		wantpoolcaount int64
	}{
		{
			name:           "test subscribe #1",
			met:            &m,
			pollInterval:   1,
			stopcount:      1,
			wantpoolcaount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ScribeMetrics(tt.met, tt.pollInterval, tt.stopcount)
			if tt.met.PollCount != tt.wantpoolcaount {
				t.Errorf("ScribeMetrics() Pollcount = %d, WantPoolCount %d", tt.met.PollCount, tt.wantpoolcaount)
			}
		})
	}
}
