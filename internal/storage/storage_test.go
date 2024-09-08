package storage

import (
	"testing"
)

func TestMemStorage_Positive_AddGauge(t *testing.T) {
	type args struct {
		name string
		item string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "positive test #1",
			args:    args{"test", "1.32445"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{
				GaugeMap: make(map[string]Gauge),
			}
			if err := storage.AddGauge(tt.args.name, tt.args.item); err != tt.wantErr {
				t.Errorf("AddGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_Negative_AddGauge(t *testing.T) {
	type args struct {
		name string
		item string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "negative test #1",
			args:    args{"test", "s"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{
				GaugeMap: make(map[string]Gauge),
			}
			if err := storage.AddGauge(tt.args.name, tt.args.item); err == tt.wantErr {
				t.Errorf("AddGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_Positive_AddCounter(t *testing.T) {
	type args struct {
		name string
		item string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "test #1",
			args:    args{"test", "1"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{
				CounterMap: make(map[string][]Counter),
			}
			if err := storage.AddCounter(tt.args.name, tt.args.item); err != tt.wantErr {
				t.Errorf("AddCounter() error = %v, dontWantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_Negative_AddCounter(t *testing.T) {
	type args struct {
		name string
		item string
	}
	tests := []struct {
		name        string
		args        args
		dontWantErr error
	}{
		{
			name:        "negative test #1",
			args:        args{"test", "1.32445"},
			dontWantErr: nil,
		},
		{
			name:        "negative test #2",
			args:        args{"test", "s"},
			dontWantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{
				CounterMap: make(map[string][]Counter),
			}
			if err := storage.AddCounter(tt.args.name, tt.args.item); err == tt.dontWantErr {
				t.Errorf("AddCounter() error = %v, dontWantErr %v", err, tt.dontWantErr)
			}
		})
	}
}
