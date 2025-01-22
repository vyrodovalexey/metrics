package scribemetrics

import (
	"github.com/vyrodovalexey/metrics/internal/model"
	"math/rand"
	"runtime"
	"time"
)

// Структура для сбора метрик
type MemMetrics struct {
	Alloc         model.Gauge   // Объем выделенной памяти
	BuckHashSys   model.Gauge   // Память, используемая для BuckHash
	Frees         model.Gauge   // Количество освобожденной памяти
	GCCPUFraction model.Gauge   // Доля CPU, используемая сборщиком мусора
	GCSys         model.Gauge   // Память, используемая сборщиком мусора
	HeapAlloc     model.Gauge   // Объем выделенной кучи
	HeapIdle      model.Gauge   // Неиспользуемая память в куче
	HeapInuse     model.Gauge   // Используемая память в куче
	HeapObjects   model.Gauge   // Количество объектов в куче
	HeapReleased  model.Gauge   // Освобожденная память в куче
	HeapSys       model.Gauge   // Общая память кучи
	LastGC        model.Gauge   // Время последней сборки мусора
	Lookups       model.Gauge   // Количество обращений к памяти
	MCacheInuse   model.Gauge   // Используемая память для кэша
	MCacheSys     model.Gauge   // Общая память для кэша
	MSpanInuse    model.Gauge   // Используемая память для MSpan
	MSpanSys      model.Gauge   // Общая память для MSpan
	Mallocs       model.Gauge   // Количество выделений памяти
	NextGC        model.Gauge   // Время до следующей сборки мусора
	NumForcedGC   model.Gauge   // Количество принудительных сборок мусора
	NumGC         model.Gauge   // Количество сборок мусора
	OtherSys      model.Gauge   // Другая системная память
	PauseTotalNs  model.Gauge   // Общее время паузы в наносекундах
	StackInuse    model.Gauge   // Используемая память стека
	StackSys      model.Gauge   // Общая память стека
	Sys           model.Gauge   // Общая системная память
	TotalAlloc    model.Gauge   // Общее количество выделенной памяти
	RandomValue   model.Gauge   // Случайное значение
	PollCount     model.Counter // Счетчик опросов
}

// Функция для сбора метрик
func UpdateMetrics(m *MemMetrics) {
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
func ScribeMetrics(m *MemMetrics, p time.Duration, stop int64) {

	for {
		// Проверяем, нужно ли остановить опрос
		if shouldStop(m.PollCount, stop) {
			return
		}
		// Собираем метрики
		UpdateMetrics(m)
		// Ждем заданный интервал времени
		<-time.After(p * time.Second)
	}
}
