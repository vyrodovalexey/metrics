package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// MemStorage Структура для хранения метрик в памяти
type MemStorage struct {
	GaugeMap   map[string]Gauge
	CounterMap map[string]Counter
}

// New Создание нового хранилища
func (m *MemStorage) New() {
	m.GaugeMap = make(map[string]Gauge)
	m.CounterMap = make(map[string]Counter)
}

// AddCounterAsString Добавление метрики Counter из строки
func (m *MemStorage) AddCounterAsString(name string, item string) error {
	counter, err := strconv.ParseInt(item, 10, 64)
	if err == nil {
		m.CounterMap[name] = m.CounterMap[name] + counter // Увеличение значения счетчика
	}
	return err
}

// AddCounter Добавление метрики Counter
func (m *MemStorage) AddCounter(name string, item Counter) {
	m.CounterMap[name] = m.CounterMap[name] + item // Увеличение значения счетчика
}

// AddGaugeAsString Добавление метрики Gauge из строки
func (m *MemStorage) AddGaugeAsString(name string, item string) error {
	gauge, err := strconv.ParseFloat(item, 64) // Парсинг строки в вещественное число
	if err == nil {
		m.GaugeMap[name] = gauge
	}
	return err
}

// AddGauge Добавление метрики Gauge
func (m *MemStorage) AddGauge(name string, item Gauge) {
	m.GaugeMap[name] = item
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
func (m *MemStorage) GetGauge(name string) (Gauge, bool) {
	res, e := m.GaugeMap[name]
	if e {
		return res, e
	}
	return 0, false
}

// GetCounter Получение метрики Counter
func (m *MemStorage) GetCounter(name string) (Counter, bool) {
	res, e := m.CounterMap[name]
	if e {
		return res, e
	}
	return 0, false
}

// Load Загрузка данных хранилища метрик из файла
func (m *MemStorage) Load(f *os.File) error {
	m.New() // Создание нового хранилища
	// Чтение содержимого файла
	byteValue, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, m)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
		}
	}
	return err
}

// SaveAsync Асинхронное сохранение данных хранилища метрик в файл
func (m *MemStorage) SaveAsync(f *os.File, interval int) {
	for {
		mst, err := json.Marshal(m)
		if err != nil {
			fmt.Println("Error move to json:", err)
		}

		// Очистка файла
		err = f.Truncate(0)
		if err != nil {
			fmt.Println("Can't truncate file error:", err)
		}
		_, err = f.Seek(0, 0) // Перемещение курсора в начало файла
		if err != nil {
			fmt.Println("Can't seek on start error:", err)
		}
		_, err = f.Write(mst) // Запись данных хранилища метрик в файл
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
		// Интервал ожидания
		<-time.After(time.Duration(interval) * time.Second)
	}
}

// Save Сохранение хранилища в файл
func (m *MemStorage) Save(f *os.File) {
	mst, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error move to json:", err)
	}
	// Очистка файла
	err = f.Truncate(0)
	if err != nil {
		fmt.Println("Can't truncate file error:", err)
	}
	_, err = f.Seek(0, 0) // Перемещение курсора в начало файла
	if err != nil {
		fmt.Println("Can't seek on start error:", err)
	}
	_, err = f.Write(mst) // Запись данных хранилища метрик в файл
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
