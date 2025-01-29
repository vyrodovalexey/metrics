package pgstorage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/vyrodovalexey/metrics/internal/model"
	"log"
)

const (
	QueryCreateGaugeTable   = "CREATE TABLE IF NOT EXISTS gauge (name TEXT NOT NULL unique,value DOUBLE PRECISION NOT NULL)"
	QueryCreateCounterTable = "CREATE TABLE IF NOT EXISTS counter (name TEXT NOT NULL unique,value INTEGER NOT NULL)"
	QuerySelectCounter      = "SELECT value FROM counter WHERE name = $1"
	QuerySelectGauge        = "SELECT value FROM gauge WHERE name = $1"
	QuerySelectAllGauge     = "SELECT name,value FROM gauge"
	QuerySelectAllCounter   = "SELECT name,value FROM counter"
	QueryInsertCounter      = "INSERT INTO counter (name,value) VALUES ($1, $2)"
	QueryInsertGauge        = "INSERT INTO gauge (name,value) VALUES ($1, $2)"
	QueryUpdateCounter      = "UPDATE counter SET value = counter.value + $2 WHERE name = $1"
	QueryUpdateGauge        = "UPDATE gauge SET value = $2 WHERE name = $1"
)

type PgStorageWithAttributes struct {
	conn    *pgx.Conn
	timeout uint
}

func (p *PgStorageWithAttributes) New(ctx context.Context, c string, timeout uint) error {
	var err error
	p.timeout = timeout
	p.conn, err = pgx.Connect(ctx, c)
	if err != nil {
		return err
	}
	_, err = p.conn.Exec(ctx, QueryCreateGaugeTable)
	if err != nil {
		return err
	}
	_, err = p.conn.Exec(ctx, QueryCreateCounterTable)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgStorageWithAttributes) Check(ctx context.Context) error {
	err := p.conn.Ping(ctx)
	return err
}

// UpdateCounter Добавление метрики Counter
func (p *PgStorageWithAttributes) UpdateCounter(ctx context.Context, name string, item model.Counter) error {
	var err error
	_, b := p.GetCounter(ctx, name)
	if b {
		_, err = p.conn.Exec(ctx, QueryUpdateCounter, name, item)
	} else {
		_, err = p.conn.Exec(ctx, QueryInsertCounter, name, item)
	}
	return err
}

// UpdateGauge Добавление метрики Gauge
func (p *PgStorageWithAttributes) UpdateGauge(ctx context.Context, name string, item model.Gauge) error {
	var err error
	_, b := p.GetGauge(ctx, name)
	if b {
		_, err = p.conn.Exec(ctx, QueryUpdateGauge, name, item)
	} else {
		_, err = p.conn.Exec(ctx, QueryInsertGauge, name, item)
	}
	return err
}

// UpdateMetric Добавление метрики в формате model. Metrics
func (p *PgStorageWithAttributes) UpdateMetric(ctx context.Context, metrics *model.Metrics) error {
	var err error
	if metrics.MType == "counter" {
		err = p.UpdateCounter(ctx, metrics.ID, *metrics.Delta)
	}
	if metrics.MType == "gauge" {
		err = p.UpdateGauge(ctx, metrics.ID, *metrics.Value)

	}
	return err
}

// GetMetric Получение метрики в формате model. Metrics
func (p *PgStorageWithAttributes) GetMetric(ctx context.Context, metrics *model.Metrics) bool {

	if metrics.MType == "counter" {
		i, b := p.GetCounter(ctx, metrics.ID)
		metrics.Delta = &i
		metrics.Value = nil
		return b
	}
	if metrics.MType == "gauge" {
		g, b := p.GetGauge(ctx, metrics.ID)
		metrics.Value = &g
		metrics.Delta = nil
		return b
	}
	return false
}

// GetAllMetricNames Получение полного списка имен метрик
func (p *PgStorageWithAttributes) GetAllMetricNames(ctx context.Context) (map[string]string, map[string]string, error) {
	var gaugerows, counterrows pgx.Rows
	var err error
	null := make(map[string]string)
	gvalues := make(map[string]string)
	cvalues := make(map[string]string)

	gaugerows, err = p.conn.Query(ctx, QuerySelectAllGauge)
	if err != nil {
		return null, null, err
	}
	defer gaugerows.Close()
	for gaugerows.Next() {
		var name string
		var v model.Gauge
		err := gaugerows.Scan(&name, &v)
		if err != nil {
			log.Fatal(err)
		}
		gvalues[name] = fmt.Sprintf("%v", v)
	}
	counterrows, err = p.conn.Query(ctx, QuerySelectAllCounter)
	if err != nil {
		return null, null, err
	}
	defer counterrows.Close()
	for counterrows.Next() {
		var name string
		var v model.Counter
		err := counterrows.Scan(&name, &v)
		if err != nil {
			return null, null, err
		}
		cvalues[name] = fmt.Sprintf("%v", v)
	}

	return gvalues, cvalues, nil
}

// GetGauge Получение метрики Gauge
func (p *PgStorageWithAttributes) GetGauge(ctx context.Context, name string) (model.Gauge, bool) {
	var value model.Gauge
	err := p.conn.QueryRow(ctx, QuerySelectGauge, name).Scan(&value)
	if err != nil {
		return 0, false
	}
	return value, true
}

// GetCounter Получение метрики Counter
func (p *PgStorageWithAttributes) GetCounter(ctx context.Context, name string) (model.Counter, bool) {
	var value model.Counter
	err := p.conn.QueryRow(ctx, QuerySelectCounter, name).Scan(&value)
	if err != nil {
		return 0, false
	}
	return value, true
}

// Load Dummy
func (p *PgStorageWithAttributes) Load(ctx context.Context, filePath string, interval uint) error {
	if ctx.Err() != nil {
		return fmt.Errorf("method is not implemented for postgresql database storage type with filepath %s and interval %d", filePath, interval)
	}
	return fmt.Errorf("method is not implemented for postgresql database storage type with filepath %s and interval %d", filePath, interval)
}

// SaveAsync Dummy
func (p *PgStorageWithAttributes) SaveAsync() error {
	return fmt.Errorf("method is not implemented for postgresql database storage type")
}

// Save Dummy
func (p *PgStorageWithAttributes) Save() error {
	return fmt.Errorf("method is not implemented for postgresql database storage type")
}

// Close Закрытие сессии
func (p *PgStorageWithAttributes) Close() {
	p.conn.Close(context.Background())
}
