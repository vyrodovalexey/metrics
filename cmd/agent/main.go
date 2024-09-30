package main

import (
	"fmt"
	"github.com/vyrodovalexey/metrics/internal/agent/config"
	"github.com/vyrodovalexey/metrics/internal/agent/sendmetrics"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const (
	maxIdleConnectionsPerHost = 10
	requestTimeout            = 30
	sendJSON                  = true
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

func httpClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnectionsPerHost,
		},
		Timeout: requestTimeout * time.Second,
	}

	return client
}

func updateMetrics(m *metrics) {
	var memStats runtime.MemStats
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
}

func shouldStop(counter int64, stop int64) bool {
	if counter < stop || stop == -1 {
		return false
	} else {
		return true
	}

}

func scribeMetrics(m *metrics, p time.Duration, stop int64) {

	for {
		if shouldStop(m.PollCount, stop) {
			return
		}
		updateMetrics(m)
		<-time.After(p * time.Second)
	}
}

func main() {

	// Создаем новый экземпляр конфигурации
	cfg := config.New()

	// Парсим настройки конфигурации
	ConfigParser(cfg)

	client := httpClient()
	m := metrics{}

	var metricSetup string
	var met sendmetrics.Metrics

	go scribeMetrics(&m, time.Duration(cfg.PoolInterval), -1)
	for {
		if m.PollCount > 0 {
			val := reflect.ValueOf(m)
			typ := reflect.TypeOf(m)
			for i := 0; i < val.NumField(); i++ {
				met.ID = typ.Field(i).Name
				// fucking setup to apply type
				switch typ.Field(i).Name {
				case "PollCount":
					metricSetup = "counter"
					met.MType = "counter"
					sint := val.Field(i).Int()
					met.Delta = &sint
				default:
					metricSetup = "gauge"
					met.MType = "gauge"
					sfloat := val.Field(i).Float()
					met.Value = &sfloat
				}
				if sendJSON {
					r := fmt.Sprintf("http://%s/update/", cfg.EndpointAddr)
					sendmetrics.SendAsJSON(client, r, &met)
				} else {
					r := fmt.Sprintf("http://%s/update/%s/%s/%v", cfg.EndpointAddr, metricSetup, typ.Field(i).Name, val.Field(i))

					sendmetrics.SendAsPlain(client, r)
				}
			}
			<-time.After(time.Duration(cfg.ReportInterval) * time.Second)
		}
	}
}
