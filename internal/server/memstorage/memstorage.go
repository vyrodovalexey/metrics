package memstorage

import (
	"encoding/json"
	"fmt"
	"github.com/vyrodovalexey/metrics/internal/model"
	"io"
	"os"
	"time"
)

// MemStorage Структура для хранения метрик в памяти
type MemStorage struct {
	GaugeMap   map[string]model.Gauge
	CounterMap map[string]model.Counter
}

// New Создание нового хранилища
func (m *MemStorage) New() {
	m.GaugeMap = make(map[string]model.Gauge)
	m.CounterMap = make(map[string]model.Counter)
}

// UpdateCounter Добавление метрики Counter
func (m *MemStorage) UpdateCounter(name string, item model.Counter, f *os.File, p bool) error {
	var err error
	m.CounterMap[name] = m.CounterMap[name] + item // Увеличение значения счетчика
	if p {
		err = m.Save(f)
	}
	return err
}

// UpdateGauge Добавление метрики Gauge
func (m *MemStorage) UpdateGauge(name string, item model.Gauge, f *os.File, p bool) error {
	var err error
	m.GaugeMap[name] = item
	if p {
		err = m.Save(f)
	}
	return err
}

// UpdateMetric Добавление метрики в формате model. Metrics
func (m *MemStorage) UpdateMetric(metrics *model.Metrics, f *os.File, p bool) error {
	var err error
	if metrics.MType == "counter" {
		err = m.UpdateCounter(metrics.ID, *metrics.Delta, f, p)
	}
	if metrics.MType == "gauge" {
		err = m.UpdateGauge(metrics.ID, *metrics.Value, f, p)

	}
	return err
}

// GetMetric Получение метрики в формате model. Metrics
func (m *MemStorage) GetMetric(metrics *model.Metrics) bool {

	if metrics.MType == "counter" {
		i, b := m.GetCounter(metrics.ID)
		metrics.Delta = &i
		metrics.Value = nil
		return b
	}
	if metrics.MType == "gauge" {
		g, b := m.GetGauge(metrics.ID)
		metrics.Value = &g
		metrics.Delta = nil
		return b
	}
	return false
}

// GetAllMetricNames Получение полного списка имен метрик
func (m *MemStorage) GetAllMetricNames() (map[string]string, map[string]string) {
	//to debug
	//names := make([]string, 0, len(storage.GaugeMap)+len(storage.CounterMap))
	gvalues := make(map[string]string, len(m.GaugeMap))
	cvalues := make(map[string]string, len(m.CounterMap))
	// Перебор карты и сбор ключей
	for name := range m.GaugeMap {

		gv, _ := m.GetGauge(name)             // Получение значения измерителя
		gvalues[name] = fmt.Sprintf("%v", gv) // Форматирование значения измерителя
	}

	for name := range m.CounterMap {

		cv, _ := m.GetCounter(name)           // Получение значения счетчика
		cvalues[name] = fmt.Sprintf("%v", cv) // Форматирование значения счетчика
	}

	return gvalues, cvalues
}

// GetGauge Получение метрики Gauge
func (m *MemStorage) GetGauge(name string) (model.Gauge, bool) {
	res, e := m.GaugeMap[name]
	if e {
		return res, e
	}
	return 0, false
}

// GetCounter Получение метрики Counter
func (m *MemStorage) GetCounter(name string) (model.Counter, bool) {
	res, e := m.CounterMap[name]
	if e {
		return res, e
	}
	return 0, false
}

// Load Загрузка данных хранилища метрик из файла
func (m *MemStorage) Load(f *os.File) error {
	var err error
	var byteValue []byte
	m.New() // Создание нового хранилища
	// Чтение содержимого файла
	byteValue, err = io.ReadAll(f)
	if err != nil {
		return err
	}
	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, m)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveAsync Асинхронное сохранение данных хранилища метрик в файл
func (m *MemStorage) SaveAsync(f *os.File, interval uint) error {
	for {
		mst, err := json.Marshal(m)
		if err != nil {
			return err
		}

		// Очистка файла
		err = f.Truncate(0)
		if err != nil {
			return err
		}
		_, err = f.Seek(0, 0) // Перемещение курсора в начало файла
		if err != nil {
			return err
		}
		_, err = f.Write(mst) // Запись данных хранилища метрик в файл
		if err != nil {
			return err
		}
		// Интервал ожидания
		<-time.After(time.Duration(interval) * time.Second)
	}
}

// Save Сохранение хранилища в файл
func (m *MemStorage) Save(f *os.File) error {
	var err error
	var mst []byte
	mst, err = json.Marshal(m)
	if err != nil {
		return err
	}
	// Очистка файла
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0) // Перемещение курсора в начало файла
	if err != nil {
		return err
	}
	_, err = f.Write(mst) // Запись данных хранилища метрик в файл
	if err != nil {
		return err
	}
	return nil
}
