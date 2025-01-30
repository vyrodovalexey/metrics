package main

import (
	"github.com/vyrodovalexey/metrics/internal/agent/config"
	"github.com/vyrodovalexey/metrics/internal/agent/scribemetrics"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	cfg := config.New()
	ConfigParser(cfg)
}

func TestScribeMetrics(t *testing.T) {

	m := scribemetrics.MemMetrics{}

	tests := []struct {
		name          string
		met           *scribemetrics.MemMetrics
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
			scribemetrics.ScribeMetrics(tt.met, tt.pollInterval, tt.stop)
			if tt.met.PollCount != tt.wantpoolcount {
				t.Errorf("ScribeMetrics() Pollcount = %d, WantPoolCount %d", tt.met.PollCount, tt.wantpoolcount)
			}
		})
	}
}
