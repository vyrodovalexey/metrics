package pgstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vyrodovalexey/metrics/internal/model"
	"go.uber.org/zap"
	"time"
)

const (
	QueryCreateGaugeTable   = "CREATE TABLE IF NOT EXISTS gauge (name TEXT NOT NULL unique,value DOUBLE PRECISION NOT NULL)"
	QueryCreateCounterTable = "CREATE TABLE IF NOT EXISTS counter (name TEXT NOT NULL unique,value bigint NOT NULL)"
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
	lg      *zap.SugaredLogger
	ctx     context.Context
}

func IsRetriableError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.ConnectionException, // General connection issue
			pgerrcode.ConnectionDoesNotExist, // Connection does not exist
			pgerrcode.ConnectionFailure,      // Failed to establish connection
			pgerrcode.QueryCanceled:          // Query canceled (e.g., due to timeout)
			return true // Retry these errors
		default:
			return false // Non-retriable errors
		}
	}
	return false
}

func (p *PgStorageWithAttributes) connectDB(ctx context.Context, c string) error {

	var err error

	for i := 0; i <= 2; i++ {
		p.lg.Infow("Connecting to database...")
		p.conn, err = pgx.Connect(ctx, c)
		if err != nil {
			if IsRetriableError(err) {
				p.lg.Infow("Database is not ready...")
			}
		} else {
			p.lg.Infow("Connected to database")
			return nil
		}

		if i == 0 {
			<-time.After(1 * time.Second)
		} else {
			<-time.After(time.Duration(i*2+1) * time.Second)
		}
	}

	p.lg.Panicw("Can't connect to database")
	return err
}

func (p *PgStorageWithAttributes) pingDB(ctx context.Context) error {

	var err error
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	for i := 0; i < 5; i++ {
		err = p.conn.Ping(timeoutCtx)
		if err != nil {
			p.lg.Infof("Database is not ready %v", err)

		} else {
			p.lg.Infow("Database ping successful")
			return nil
		}

		if i == 0 {
			<-time.After(1 * time.Second)
		} else {
			<-time.After(time.Duration(i*2+1) * time.Second)
		}
	}

	p.lg.Infow("Can't connect to database")
	return err
}

func (p *PgStorageWithAttributes) execDB(ctx context.Context, query string, args ...any) error {

	var err error
	err = p.pingDB(ctx)
	if err != nil {
		return err
	}
	for i := 0; i < 5; i++ {

		_, err = p.conn.Exec(ctx, query, args...)

		if err != nil {
			p.lg.Infof("Database is not ready to execute query %v", err)
		} else {
			p.lg.Infof("Query %s executed", query)
			return nil
		}

		if i == 0 {
			<-time.After(1 * time.Second)
		} else {
			<-time.After(time.Duration(i*2+1) * time.Second)
		}
	}

	p.lg.Infow("Can't connect to database to execute query")
	return err
}

func (p *PgStorageWithAttributes) New(c string, timeout uint, log *zap.SugaredLogger) error {
	var err error
	p.ctx = context.Background()
	p.timeout = timeout
	p.lg = log
	err = p.connectDB(p.ctx, c)
	if err != nil {
		return err
	}
	err = p.execDB(p.ctx, QueryCreateGaugeTable)
	if err != nil {
		return err
	}
	err = p.execDB(p.ctx, QueryCreateCounterTable)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgStorageWithAttributes) Check() error {
	err := p.pingDB(p.ctx)
	return err
}

// UpdateCounter Добавление метрики Counter
func (p *PgStorageWithAttributes) UpdateCounter(name string, item model.Counter) error {
	var err error
	_, b := p.GetCounter(name)
	if b {
		err = p.execDB(p.ctx, QueryUpdateCounter, name, item)
	} else {
		err = p.execDB(p.ctx, QueryInsertCounter, name, item)
	}
	return err
}

// UpdateGauge Добавление метрики Gauge
func (p *PgStorageWithAttributes) UpdateGauge(name string, item model.Gauge) error {
	var err error
	_, b := p.GetGauge(name)
	if b {
		err = p.execDB(p.ctx, QueryUpdateGauge, name, item)
	} else {
		err = p.execDB(p.ctx, QueryInsertGauge, name, item)
	}
	return err
}

// UpdateMetric Добавление метрики в формате model. Metrics
func (p *PgStorageWithAttributes) UpdateMetric(metrics *model.Metrics) error {
	var err error
	if metrics.MType == "counter" {
		err = p.UpdateCounter(metrics.ID, *metrics.Delta)
	}
	if metrics.MType == "gauge" {
		err = p.UpdateGauge(metrics.ID, *metrics.Value)

	}
	return err
}

// GetMetric Получение метрики в формате model. Metrics
func (p *PgStorageWithAttributes) GetMetric(metrics *model.Metrics) bool {

	if metrics.MType == "counter" {
		i, b := p.GetCounter(metrics.ID)
		metrics.Delta = &i
		metrics.Value = nil
		return b
	}
	if metrics.MType == "gauge" {
		g, b := p.GetGauge(metrics.ID)
		metrics.Value = &g
		metrics.Delta = nil
		return b
	}
	return false
}

// GetAllMetricNames Получение полного списка имен метрик
func (p *PgStorageWithAttributes) GetAllMetricNames() (map[string]string, map[string]string, error) {
	var gaugerows, counterrows pgx.Rows
	var err error
	null := make(map[string]string)
	gvalues := make(map[string]string)
	cvalues := make(map[string]string)
	err = p.pingDB(p.ctx)
	if err != nil {
		return null, null, err
	}

	gaugerows, err = p.conn.Query(p.ctx, QuerySelectAllGauge)
	if err != nil {
		return null, null, err
	}
	defer gaugerows.Close()
	for gaugerows.Next() {
		var name string
		var v model.Gauge
		err := gaugerows.Scan(&name, &v)
		if err != nil {
			p.lg.Infow("gauge scan", "error", err)
			return null, null, err
		}
		gvalues[name] = fmt.Sprintf("%v", v)
	}
	counterrows, err = p.conn.Query(p.ctx, QuerySelectAllCounter)
	if err != nil {
		return null, null, err
	}
	defer counterrows.Close()
	for counterrows.Next() {
		var name string
		var v model.Counter
		err := counterrows.Scan(&name, &v)
		if err != nil {
			p.lg.Infow("counter scan", "error", err)
			return null, null, err
		}
		cvalues[name] = fmt.Sprintf("%v", v)
	}

	return gvalues, cvalues, nil
}

// GetGauge Получение метрики Gauge
func (p *PgStorageWithAttributes) GetGauge(name string) (model.Gauge, bool) {
	var value model.Gauge
	var err error
	err = p.pingDB(p.ctx)
	if err != nil {
		return 0, false
	}
	err = p.conn.QueryRow(p.ctx, QuerySelectGauge, name).Scan(&value)
	if err != nil {
		return 0, false
	}
	return value, true
}

// GetCounter Получение метрики Counter
func (p *PgStorageWithAttributes) GetCounter(name string) (model.Counter, bool) {
	var value model.Counter
	var err error
	err = p.pingDB(p.ctx)
	if err != nil {
		return 0, false
	}
	err = p.conn.QueryRow(p.ctx, QuerySelectCounter, name).Scan(&value)
	if err != nil {
		return 0, false
	}
	return value, true
}

// Load Dummy
func (p *PgStorageWithAttributes) Load(filePath string, interval uint, log *zap.SugaredLogger) error {
	if p.ctx.Err() != nil {
		return fmt.Errorf("operation interapted for method which implemented for postgresql database storage type with filepath %s and interval %d", filePath, interval)
	}
	return fmt.Errorf("method is not implemented for postgresql database storage type with filepath %s and interval %d", filePath, interval)
}

// SaveAsync Dummy
func (p *PgStorageWithAttributes) SaveAsync() {
	p.lg.Infow("Method is not implemented for postgresql database storage type")
}

// Save Dummy
func (p *PgStorageWithAttributes) Save() error {
	return fmt.Errorf("method is not implemented for postgresql database storage type")
}

// Close Закрытие сессии
func (p *PgStorageWithAttributes) Close() {
	p.conn.Close(context.Background())
}
