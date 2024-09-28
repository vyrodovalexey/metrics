package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const (
	serverAddr            = "localhost:8080"
	defaultReportInterval = 10
	defaultPoolInterval   = 2
	sendjson              = true
)

type Config struct {
	EndpointAddr   string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PoolInterval   int    `env:"POLL_INTERVAL"`
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

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

func SendMetricPlain(cl http.Client, url string) {
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

func SendMetricJSON(cl http.Client, url string, m *Metrics) {
	jm, _ := json.Marshal(*m)
	req, errr := http.NewRequest("POST", url, bytes.NewBuffer(jm))
	if errr != nil {
		log.Fatal(errr)
	}
	req.Header.Set("Content-Type", "application/json")

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

	var cfg Config
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	if len(cfg.EndpointAddr) == 0 {
		flag.StringVar(&cfg.EndpointAddr, "a", serverAddr, "input ip:port or host:port of metrics server")
	}
	if cfg.ReportInterval < 1 {
		flag.IntVar(&cfg.ReportInterval, "r", defaultReportInterval, "seconds delay interval to send metrics to metrics server")
	}
	if cfg.PoolInterval < 1 {
		flag.IntVar(&cfg.PoolInterval, "p", defaultPoolInterval, "seconds delay between scribing metrics from host")
	}
	flag.Parse()
	client := &http.Client{}
	m := metrics{}
	// variable for setup
	var metrict string

	var met Metrics

	go ScribeMetrics(&m, time.Duration(cfg.PoolInterval), -1)
	for {

		if m.PollCount > 0 {
			val := reflect.ValueOf(m)
			typ := reflect.TypeOf(m)
			for i := 0; i < val.NumField(); i++ {
				met.ID = typ.Field(i).Name
				// fucking setup to apply type
				switch typ.Field(i).Name {
				case "PollCount":
					metrict = "counter"
					met.MType = "counter"
					sint := val.Field(i).Int()
					met.Delta = &sint
				default:
					metrict = "gauge"
					met.MType = "gauge"
					sfloat := val.Field(i).Float()
					met.Value = &sfloat
				}
				if sendjson {
					r := fmt.Sprintf("http://%s/update/", cfg.EndpointAddr)
					SendMetricJSON(*client, r, &met)
				} else {
					r := fmt.Sprintf("http://%s/update/%s/%s/%v", cfg.EndpointAddr, metrict, typ.Field(i).Name, val.Field(i))

					SendMetricPlain(*client, r)
				}
				time.Sleep(time.Duration(cfg.ReportInterval) * time.Second)
			}
		}
	}
}
