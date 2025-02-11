package model

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

const (
	CounterString = "counter"
	GaugeString   = "gauge"
)

type Gauge = float64
type Counter = int64

// Metrics Структура для обмена метрик
type Metrics struct {
	ID    string   `json:"id" binding:"required"`                       // Имя метрики
	MType string   `json:"type" binding:"required,oneof=counter gauge"` // Тип метрики (gauge или counter)
	Delta *int64   `json:"delta,omitempty"`                             // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"`                             // Значение метрики в случае передачи gauge
}

type MetricsBatch []Metrics

// URLPathToMetric парсит параметры URL и заполняет структуру Metrics.
func (mod *Metrics) URLPathToMetric(mtype string, key string, value string) error {
	var err error
	var gauge float64
	var counter int64
	switch mtype {
	case GaugeString:
		// Если тип метрики - gauge, устанавливаем тип и имя метрики
		mod.MType = GaugeString
		mod.ID = key
		// Если значение value не пустое, парсим его как float64 и устанавливаем в поле Value
		if len(value) > 0 {
			gauge, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			mod.Value = &gauge
		}
	case CounterString:
		// Если тип метрики - counter, устанавливаем тип и имя метрики
		mod.MType = CounterString
		mod.ID = key
		// Если значение value не пустое, парсим его как int64 и устанавливаем в поле Delta
		if len(value) > 0 {
			counter, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			mod.Delta = &counter
		}

	default:
		// Если тип метрики не gauge и не counter, возвращаем ошибку
		err = fmt.Errorf("unknown metric type: %s", mod.MType)
	}
	return err
}

// BodyToMetric декодирует JSON из тела запроса в структуру Metrics.
func (mod *Metrics) BodyToMetric(body io.ReadCloser) error {
	var err error
	// Декодируем JSON из тела запроса в структуру Metrics
	err = json.NewDecoder(body).Decode(mod)
	if err != nil {
		return err
	}
	// Проверяем тип метрики

	switch mod.MType {
	case GaugeString:
		err = nil

	case CounterString:
		err = nil

	default:
		// Если тип метрики не gauge и не counter, возвращаем ошибку
		err = fmt.Errorf("unknown metric type: %s", mod.MType)
	}

	return err
}

func (batch *MetricsBatch) BodyToMetricBatch(body io.ReadCloser) error {
	var err error
	// Декодируем JSON из тела запроса в структуру Metrics
	err = json.NewDecoder(body).Decode(batch)
	if err != nil {
		return err
	}
	for i := range *batch {
		// Проверяем тип метрики
		switch (*batch)[i].MType {
		case GaugeString:
			err = nil

		case CounterString:
			err = nil

		default:
			// Если тип метрики не gauge и не counter, возвращаем ошибку
			err = fmt.Errorf("unknown metric type: %s, for metric id %s", (*batch)[i].MType, (*batch)[i].ID)
		}
	}
	return err
}

func (mod *Metrics) String() string {
	var res string
	switch mod.MType {
	case GaugeString:
		res = fmt.Sprintf("%.3f", *mod.Value)
	case CounterString:
		res = fmt.Sprintf("%d", *mod.Delta)
	}
	return res

}
