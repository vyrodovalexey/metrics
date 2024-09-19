package main

import (
	"fmt"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const (
	pollInterval   = 2
	reportInterval = 10
	address        = "localhost:8080"
	stopCount      = -1
)

type metrics struct {
	Alloc         storage.Gauge
	BuckHashSys   storage.Gauge
	Frees         storage.Gauge
	GCCPUFraction storage.Gauge
	GCSys         storage.Gauge
	HeapAlloc     storage.Gauge
	HeapIdle      storage.Gauge
	HeapInuse     storage.Gauge
	HeapObjects   storage.Gauge
	HeapReleased  storage.Gauge
	HeapSys       storage.Gauge
	LastGC        storage.Gauge
	Lookups       storage.Gauge
	MCacheInuse   storage.Gauge
	MCacheSys     storage.Gauge
	MSpanInuse    storage.Gauge
	MSpanSys      storage.Gauge
	Mallocs       storage.Gauge
	NextGC        storage.Gauge
	NumForcedGC   storage.Gauge
	NumGC         storage.Gauge
	OtherSys      storage.Gauge
	PauseTotalNs  storage.Gauge
	StackInuse    storage.Gauge
	StackSys      storage.Gauge
	Sys           storage.Gauge
	TotalAlloc    storage.Gauge
	RandomValue   storage.Gauge
	PollCount     storage.Counter
}

func SendMetric(cl http.Client, url string) {
	req, errr := http.NewRequest("POST", url, nil)
	if errr != nil {
		log.Fatal(errr)
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := cl.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(time.Now(), " ", url, " ", resp.StatusCode)
	defer resp.Body.Close()
}

func ScribeMetrics(m *metrics, p time.Duration, stop int64) {
	var memStats runtime.MemStats

	for {
		if m.PollCount >= stop && stop != -1 {
			return
		} else {
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
			time.Sleep(p * time.Second)
		}
	}
}

func main() {
	client := &http.Client{}
	m := metrics{}
	// variable for setup
	var metrict string
	go ScribeMetrics(&m, pollInterval, stopCount)
	for {
		time.Sleep(reportInterval * time.Second)
		if m.PollCount > 0 {
			val := reflect.ValueOf(m)
			typ := reflect.TypeOf(m)
			for i := 0; i < val.NumField(); i++ {
				// fucking setup to apply type
				switch typ.Field(i).Name {
				case "PollCount":
					metrict = "counter"
				default:
					metrict = "gauge"
				}
				r := fmt.Sprintf("http://%s/update/%s/%s/%v", address, metrict, typ.Field(i).Name, val.Field(i))
				SendMetric(*client, r)
			}
		}
	}
}
