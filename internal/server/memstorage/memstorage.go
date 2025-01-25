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

type MemStorageWithAttributes struct {
	mst      MemStorage
	f        *os.File
	p        bool
	interval uint
}

// New Создание нового хранилища
func (m *MemStorageWithAttributes) NewMemStorage(filePath string, interval uint) error {
	var err error
	m.mst.GaugeMap = make(map[string]model.Gauge)
	m.mst.CounterMap = make(map[string]model.Counter)
	// Открываем или создаем файл для хранения
	m.f, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	m.interval = interval
	if interval > 0 {
		m.p = true
	} else {
		m.p = false
	}
	return nil
}

// UpdateCounter Добавление метрики Counter
func (m *MemStorageWithAttributes) UpdateCounter(name string, item model.Counter) error {
	var err error
	m.mst.CounterMap[name] = m.mst.CounterMap[name] + item // Увеличение значения счетчика
	if m.p {
		err = m.Save()
	}
	return err
}

// UpdateGauge Добавление метрики Gauge
func (m *MemStorageWithAttributes) UpdateGauge(name string, item model.Gauge) error {
	var err error
	m.mst.GaugeMap[name] = item
	if m.p {
		err = m.Save()
	}
	return err
}

// UpdateMetric Добавление метрики в формате model. Metrics
func (m *MemStorageWithAttributes) UpdateMetric(metrics *model.Metrics) error {
	var err error
	if metrics.MType == "counter" {
		err = m.UpdateCounter(metrics.ID, *metrics.Delta)
	}
	if metrics.MType == "gauge" {
		err = m.UpdateGauge(metrics.ID, *metrics.Value)

	}
	if m.p {
		err = m.Save()
	}
	return err
}

// GetMetric Получение метрики в формате model. Metrics
func (m *MemStorageWithAttributes) GetMetric(metrics *model.Metrics) bool {

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
func (m *MemStorageWithAttributes) GetAllMetricNames() (map[string]string, map[string]string) {
	//to debug
	//names := make([]string, 0, len(storage.GaugeMap)+len(storage.CounterMap))
	gvalues := make(map[string]string, len(m.mst.GaugeMap))
	cvalues := make(map[string]string, len(m.mst.CounterMap))
	// Перебор карты и сбор ключей
	for name := range m.mst.GaugeMap {

		gv, _ := m.GetGauge(name)             // Получение значения измерителя
		gvalues[name] = fmt.Sprintf("%v", gv) // Форматирование значения измерителя
	}

	for name := range m.mst.CounterMap {

		cv, _ := m.GetCounter(name)           // Получение значения счетчика
		cvalues[name] = fmt.Sprintf("%v", cv) // Форматирование значения счетчика
	}

	return gvalues, cvalues
}

// GetGauge Получение метрики Gauge
func (m *MemStorageWithAttributes) GetGauge(name string) (model.Gauge, bool) {
	res, e := m.mst.GaugeMap[name]
	if e {
		return res, e
	}
	return 0, false
}

// GetCounter Получение метрики Counter
func (m *MemStorageWithAttributes) GetCounter(name string) (model.Counter, bool) {
	res, e := m.mst.CounterMap[name]
	if e {
		return res, e
	}
	return 0, false
}

// Load Загрузка данных хранилища метрик из файла
func (m *MemStorageWithAttributes) LoadMemStorage(filePath string, interval uint) error {
	var err error
	var byteValue []byte
	m.NewMemStorage(filePath, interval) // Создание нового хранилища
	// Чтение содержимого файла
	byteValue, err = io.ReadAll(m.f)
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
func (m *MemStorageWithAttributes) SaveAsync() error {
	for {
		mst, err := json.Marshal(m.mst)
		if err != nil {
			return err
		}

		// Очистка файла
		err = m.f.Truncate(0)
		if err != nil {
			return err
		}
		_, err = m.f.Seek(0, 0) // Перемещение курсора в начало файла
		if err != nil {
			return err
		}
		_, err = m.f.Write(mst) // Запись данных хранилища метрик в файл
		if err != nil {
			return err
		}
		// Интервал ожидания
		<-time.After(time.Duration(m.interval) * time.Second)
	}
}

// Save Сохранение хранилища в файл
func (m *MemStorageWithAttributes) Save() error {
	var err error
	var mst []byte
	mst, err = json.Marshal(m.mst)
	if err != nil {
		return err
	}
	// Очистка файла
	err = m.f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = m.f.Seek(0, 0) // Перемещение курсора в начало файла
	if err != nil {
		return err
	}
	_, err = m.f.Write(mst) // Запись данных хранилища метрик в файл
	if err != nil {
		return err
	}
	return nil
}

func (m *MemStorageWithAttributes) Close() {
	//	defer conn.Close(ctx)
	defer m.f.Close()
}
