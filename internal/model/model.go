package model

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// Metrics Структура для обмена метрик
type Metrics struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // Тип метрики (gauge или counter)
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
}

// URLPathToMetric парсит параметры URL и заполняет структуру Metrics.
func (mod *Metrics) URLPathToMetric(mtype string, key string, value string) error {
	var err error
	var gauge float64
	var counter int64
	switch mtype {
	case "gauge":
		// Если тип метрики - gauge, устанавливаем тип и имя метрики
		mod.MType = "gauge"
		mod.ID = key
		// Если значение value не пустое, парсим его как float64 и устанавливаем в поле Value
		if len(value) > 0 {
			gauge, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			mod.Value = &gauge
		}
	case "counter":
		// Если тип метрики - counter, устанавливаем тип и имя метрики
		mod.MType = "counter"
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
func (mod *Metrics) BodyToMetric(body *io.ReadCloser) error {
	var err error
	// Декодируем JSON из тела запроса в структуру Metrics
	err = json.NewDecoder(*body).Decode(mod)
	if err != nil {
		return err
	}
	// Проверяем тип метрики
	switch mod.MType {
	case "gauge":
		err = nil

	case "counter":
		err = nil

	default:
		// Если тип метрики не gauge и не counter, возвращаем ошибку
		err = fmt.Errorf("unknown metric type: %s", mod.MType)

	}
	return err
}
