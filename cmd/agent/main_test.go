package main

import (
	"context"
	"fmt"
	"github.com/vyrodovalexey/metrics/internal/agent/config"
	"github.com/vyrodovalexey/metrics/internal/agent/scribemetrics"
	"github.com/vyrodovalexey/metrics/internal/agent/sendmetrics"
	"github.com/vyrodovalexey/metrics/internal/model"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"syscall"
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

func TestSendRequests(t *testing.T) {
	var err error
	s := 2.12
	m := model.Metrics{ID: "test", MType: "gauge", Value: &s}
	b := model.MetricsBatch{}
	b = append(b, m)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	client := httpClient()
	err = sendmetrics.SendAsJSON(client, ts.URL, &m)
	if err != nil {
		t.Fatalf("Failed to send JSON request: %v", err)
	}
	err = sendmetrics.SendAsBatchJSON(client, ts.URL, &b)
	if err != nil {
		t.Fatalf("Failed to send JSONBatch request: %v", err)
	}
	err = sendmetrics.SendAsPlain(client, fmt.Sprintf("%s/update/gauge/test/2.12", ts.URL))
	if err != nil {
		t.Fatalf("Failed to send JSON request: %v", err)
	}
}

func TestSendRequestsWithServerUnavailable(t *testing.T) {
	var err error

	s := 2.12
	m := model.Metrics{ID: "test", MType: "gauge", Value: &s}
	b := model.MetricsBatch{}
	b = append(b, m)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)

	}))
	defer ts.Close()
	var dialAttempt int32
	customTransport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Fail first connection attempt
			if atomic.AddInt32(&dialAttempt, 1) < 3 {
				return nil, &net.OpError{
					Op:  "dial",
					Net: network,
					Err: syscall.ECONNREFUSED,
				}
			}
			// Subsequent attempts use default dialer
			return (&net.Dialer{}).DialContext(ctx, network, addr)
		},
	}

	client := &http.Client{
		Transport: customTransport,
	}

	err = sendmetrics.SendAsJSON(client, ts.URL, &m)
	if err != nil {
		t.Fatalf("Failed to send JSON request: %v", err)
	}
	err = sendmetrics.SendAsBatchJSON(client, ts.URL, &b)
	if err != nil {
		t.Fatalf("Failed to send JSONBatch request: %v", err)
	}
	err = sendmetrics.SendAsPlain(client, fmt.Sprintf("%s/update/gauge/test/2.12", ts.URL))
	if err != nil {
		t.Fatalf("Failed to send JSON request: %v", err)
	}
}
