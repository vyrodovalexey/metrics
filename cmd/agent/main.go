package main

import (
	"fmt"
	"github.com/vyrodovalexey/metrics/internal/agent/config"
	"github.com/vyrodovalexey/metrics/internal/agent/scribemetrics"
	"github.com/vyrodovalexey/metrics/internal/agent/sendmetrics"
	"github.com/vyrodovalexey/metrics/internal/model"
	"log"
	"net/http"
	"reflect"
	"time"
)

const (
	maxIdleConnectionsPerHost = 10    // Максимальное количество неактивных соединений на хост
	requestTimeout            = 30    // Таймаут запроса в секундах
	sendJSON                  = false // Флаг для отправки данных в формате JSON
)

// Функция для создания HTTP клиента
func httpClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnectionsPerHost, // Устанавливаем максимальное количество неактивных соединений
		},
		Timeout: requestTimeout * time.Second, // Устанавливаем таймаут запроса
	}

	return client
}

func main() {

	// Инициализируем новый экземпляр конфигурации и
	// парсим настройки конфигурации
	cfg := config.New()
	ConfigParser(cfg)

	client := httpClient()

	//ctx := context.Background()

	// Инициализируем структуру метрик
	m := scribemetrics.MemMetrics{}

	// Костыль - переменная для хранения типа метрики
	var metricSetup string

	// Инициализируем структуру для метрик
	var met model.Metrics

	// Запускаем горутину для сбора метрик
	go scribemetrics.ScribeMetrics(&m, time.Duration(cfg.PoolInterval), -1)
	for {
		if m.PollCount > 0 {
			var err error
			val := reflect.ValueOf(m)
			typ := reflect.TypeOf(m)
			// Инициализируем структуру для метрик в виде batch
			batch := make(model.MetricsBatch, val.NumField())
			for i := 0; i < val.NumField(); i++ {
				met.ID = typ.Field(i).Name
				// Костыль - Настройка типа метрики
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
				// Выбираем отправку в формате BatchJSON, JSON или plaintext
				switch {
				case cfg.BatchMode:
					batch[i] = met
				case sendJSON:
					r := fmt.Sprintf("http://%s/update/", cfg.EndpointAddr)
					err = sendmetrics.SendAsJSON(client, r, &met)
				default:
					r := fmt.Sprintf("http://%s/update/%s/%s/%v", cfg.EndpointAddr, metricSetup, typ.Field(i).Name, val.Field(i))
					err = sendmetrics.SendAsPlain(client, r)
				}
				if err != nil {
					log.Println(err)
				}
			}
			if cfg.BatchMode {
				r := fmt.Sprintf("http://%s/updates/", cfg.EndpointAddr)
				err = sendmetrics.SendAsBatchJSON(client, r, &batch)
			}
			if err != nil {
				log.Println(err)
			}
		}
		<-time.After(time.Duration(cfg.ReportInterval) * time.Second)
	}
}
