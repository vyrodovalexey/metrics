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
			storage := &MemStorage{}
			storage.Init()
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
			storage := &MemStorage{}
			storage.Init()
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
			storage := &MemStorage{}
			storage.Init()
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
			storage := &MemStorage{}
			storage.Init()
			if err := storage.AddCounter(tt.args.name, tt.args.item); err == tt.dontWantErr {
				t.Errorf("AddCounter() error = %v, dontWantErr %v", err, tt.dontWantErr)
			}
		})
	}
}

func TestMemStorage_Positive_GetAll(t *testing.T) {

	tests := []struct {
		name           string
		countermapsize int
		gaugemapsize   int
	}{
		{
			name:           "positive get test #1",
			countermapsize: 1,
			gaugemapsize:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{}
			storage.Init()
			storage.AddCounter("test", "12")
			storage.AddGauge("test", "12.1")
			gres, cres := storage.GetAllMetricNames()
			if len(cres) < tt.countermapsize {
				t.Errorf("Size of ConterMap is %d and it less then required size %d", len(cres), tt.countermapsize)
			}
			if len(gres) < tt.gaugemapsize {
				t.Errorf("Size of GaugeMap is %d and it less then required size %d", len(gres), tt.gaugemapsize)
			}
		})
	}
}
