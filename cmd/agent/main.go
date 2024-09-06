package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

type gauge = float64
type counter = int64

type metrics struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	RandomValue   gauge
	PollCount     counter
}

var metrictype string

func send_metric(cl http.Client, url string) {
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "text/plain")
	_, err = cl.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
}

func main() {
	client := &http.Client{}

	var memStats runtime.MemStats
	m := metrics{}
	for {
		runtime.ReadMemStats(&memStats)
		m.Alloc = float64(memStats.Alloc)
		m.BuckHashSys = float64(memStats.BuckHashSys)
		m.GCCPUFraction = memStats.GCCPUFraction
		m.Frees = float64(memStats.Frees)
		m.GCSys = float64(memStats.GCSys)
		m.HeapAlloc = float64(memStats.HeapAlloc)
		m.HeapIdle = float64(memStats.HeapIdle)
		m.HeapInuse = float64(memStats.HeapInuse)
		m.HeapObjects = float64(memStats.HeapObjects)
		m.HeapReleased = float64(memStats.HeapReleased)
		m.HeapSys = float64(memStats.HeapSys)
		m.LastGC = float64(memStats.LastGC)
		m.Lookups = float64(memStats.Lookups)
		m.MCacheInuse = float64(memStats.MCacheInuse)
		m.MCacheSys = float64(memStats.MCacheSys)
		m.MSpanInuse = float64(memStats.MSpanInuse)
		m.MSpanSys = float64(memStats.MSpanSys)
		m.Mallocs = float64(memStats.Mallocs)
		m.NextGC = float64(memStats.NextGC)
		m.NumForcedGC = float64(memStats.NumForcedGC)
		m.NumGC = float64(memStats.NumGC)
		m.OtherSys = float64(memStats.OtherSys)
		m.PauseTotalNs = float64(memStats.PauseTotalNs)
		m.StackInuse = float64(memStats.StackInuse)
		m.StackSys = float64(memStats.StackSys)
		m.Sys = float64(memStats.Sys)
		m.TotalAlloc = float64(memStats.TotalAlloc)
		m.RandomValue = rand.Float64()
		m.PollCount += 1

		if m.PollCount%5 == 0 {
			val := reflect.ValueOf(m)
			typ := reflect.TypeOf(m)

			for i := 0; i < val.NumField(); i++ {
				if typ.Field(i).Name == "PollCount" {
					metrictype = "counter"
				} else {
					metrictype = "gauge"
				}
				r := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", metrictype, typ.Field(i).Name, val.Field(i))
				send_metric(*client, r)
			}
		}
		time.Sleep(2 * time.Second)
	}
}
