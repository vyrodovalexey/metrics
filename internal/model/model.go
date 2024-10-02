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

func (mod *Metrics) URLPathToMetric(mtype string, key string, value string) error {
	var err error
	var gauge float64
	var counter int64
	switch mtype {
	case "gauge":

		mod.MType = "gauge"
		mod.ID = key
		if len(value) > 0 {
			gauge, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			mod.Value = &gauge
		}
	case "counter":

		mod.MType = "counter"
		mod.ID = key
		if len(value) > 0 {
			counter, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			mod.Delta = &counter
		}

	default:
		err = fmt.Errorf("unknown metric type: %s", mod.MType)
	}
	return err
}

func (mod *Metrics) BodyToMetric(body *io.ReadCloser) error {
	var err error
	err = json.NewDecoder(*body).Decode(mod)
	if err != nil {
		return err
	}
	switch mod.MType {
	case "gauge":
		err = nil

	case "counter":
		err = nil

	default:
		err = fmt.Errorf("unknown metric type: %s", mod.MType)

	}
	return err
}
