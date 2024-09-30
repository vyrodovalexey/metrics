package main

import (
	"fmt"
	"github.com/vyrodovalexey/metrics/internal/agent/config"
	"github.com/vyrodovalexey/metrics/internal/agent/sendmetrics"
	"github.com/vyrodovalexey/metrics/internal/model"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const (
	maxIdleConnectionsPerHost = 10   // Максимальное количество неактивных соединений на хост
	requestTimeout            = 30   // Таймаут запроса в секундах
	sendJSON                  = true // Флаг для отправки данных в формате JSON
)

// Структура для сбора метрик
type metrics struct {
	Alloc         storage.Gauge   // Объем выделенной памяти
	BuckHashSys   storage.Gauge   // Память, используемая для BuckHash
	Frees         storage.Gauge   // Количество освобожденной памяти
	GCCPUFraction storage.Gauge   // Доля CPU, используемая сборщиком мусора
	GCSys         storage.Gauge   // Память, используемая сборщиком мусора
	HeapAlloc     storage.Gauge   // Объем выделенной кучи
	HeapIdle      storage.Gauge   // Неиспользуемая память в куче
	HeapInuse     storage.Gauge   // Используемая память в куче
	HeapObjects   storage.Gauge   // Количество объектов в куче
	HeapReleased  storage.Gauge   // Освобожденная память в куче
	HeapSys       storage.Gauge   // Общая память кучи
	LastGC        storage.Gauge   // Время последней сборки мусора
	Lookups       storage.Gauge   // Количество обращений к памяти
	MCacheInuse   storage.Gauge   // Используемая память для кэша
	MCacheSys     storage.Gauge   // Общая память для кэша
	MSpanInuse    storage.Gauge   // Используемая память для MSpan
	MSpanSys      storage.Gauge   // Общая память для MSpan
	Mallocs       storage.Gauge   // Количество выделений памяти
	NextGC        storage.Gauge   // Время до следующей сборки мусора
	NumForcedGC   storage.Gauge   // Количество принудительных сборок мусора
	NumGC         storage.Gauge   // Количество сборок мусора
	OtherSys      storage.Gauge   // Другая системная память
	PauseTotalNs  storage.Gauge   // Общее время паузы в наносекундах
	StackInuse    storage.Gauge   // Используемая память стека
	StackSys      storage.Gauge   // Общая память стека
	Sys           storage.Gauge   // Общая системная память
	TotalAlloc    storage.Gauge   // Общее количество выделенной памяти
	RandomValue   storage.Gauge   // Случайное значение
	PollCount     storage.Counter // Счетчик опросов
}

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

// Функция для сбора метрик
func updateMetrics(m *metrics) {
	// Структура для хранения статистики памяти
	var memStats runtime.MemStats
	// Читаем статистику памяти
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

// Функция для проверки, нужно ли остановить опрос
func shouldStop(counter int64, stop int64) bool {
	if counter < stop || stop == -1 { // Если счетчик меньше значения остановки или значение остановки равно -1
		return false // Не останавливаем
	} else {
		return true // Останавливаем
	}
}

// Функция для записи метрик
func scribeMetrics(m *metrics, p time.Duration, stop int64) {

	for {
		// Проверяем, нужно ли остановить опрос
		if shouldStop(m.PollCount, stop) {
			return
		}
		// Собираем метрики
		updateMetrics(m)
		// Ждем заданный интервал времени
		<-time.After(p * time.Second)
	}
}

func main() {

	// Инициализируем новый экземпляр конфигурации и
	// парсим настройки конфигурации
	cfg := config.New()
	ConfigParser(cfg)

	client := httpClient()
	// Инициализируем структуру метрик
	m := metrics{}
	// Костыль - переменная для хранения типа метрики
	var metricSetup string

	// Инициализируем структуру для метрик
	var met model.Metrics

	// Запускаем горутину для сбора метрик
	go scribeMetrics(&m, time.Duration(cfg.PoolInterval), -1)
	for {
		if m.PollCount > 0 {
			val := reflect.ValueOf(m)
			typ := reflect.TypeOf(m)
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
				// Выбираем отправку в формате JSON или plaintext
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
