package main

import (
	"github.com/vyrodovalexey/metrics/internal/model"
	"github.com/vyrodovalexey/metrics/internal/server/config"
	"github.com/vyrodovalexey/metrics/internal/server/logging"
	"github.com/vyrodovalexey/metrics/internal/server/memstorage"
	"github.com/vyrodovalexey/metrics/internal/server/routing"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRequest(t *testing.T) {
	gauge := make(map[string]model.Gauge)
	counter := make(map[string]model.Counter)
	mst := memstorage.MemStorage{GaugeMap: gauge, CounterMap: counter}

	sugar := logging.NewLogging(zap.InfoLevel)
	// Создаем новый экземпляр конфигурации
	cfg := config.New()
	cfg.FileStoragePath = "/tmp/test.json"
	// Парсим настройки конфигурации
	ConfigParser(cfg)
	var awf bool
	if cfg.StoreInterval > 0 {
		awf = false
	} else {
		awf = true
	}

	file, _ := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	router := routing.SetupRouter(sugar)
	routing.ConfigureRouting(router, &mst, file, awf)
	router.LoadHTMLGlob("../../templates/*")
	tests := []struct {
		name           string
		st             *memstorage.MemStorage
		method         string
		url            string
		mimetype       string
		expectedStatus int
	}{
		{
			name:           "Valid Update gauge",
			st:             &mst,
			method:         http.MethodPost,
			url:            "/update/gauge/test/1.4545",
			mimetype:       "text/plain",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid Update counter",
			st:             &mst,
			method:         http.MethodPost,
			url:            "/update/counter/test/1",
			mimetype:       "text/plain",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Method",
			st:             &mst,
			method:         http.MethodGet,
			url:            "/update/gauge/test/1.12",
			mimetype:       "text/plain",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid Gauge",
			st:             &mst,
			method:         http.MethodPost,
			url:            "/update/gauge/test/test",
			mimetype:       "text/plain",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Counter",
			st:             &mst,
			method:         http.MethodPost,
			url:            "/update/counter/test/1.12",
			mimetype:       "text/plain",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Get counter",
			st:             &mst,
			method:         http.MethodGet,
			url:            "/value/counter/test",
			mimetype:       "text/plain",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get gauge",
			st:             &mst,
			method:         http.MethodGet,
			url:            "/value/gauge/test",
			mimetype:       "text/plain",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Get gauge",
			st:             &mst,
			method:         http.MethodGet,
			url:            "/value/gauge/unavailable",
			mimetype:       "text/plain",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Get /",
			st:             &mst,
			method:         http.MethodGet,
			url:            "/",
			mimetype:       "text/html",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			req.Header.Add("Content-Type", tt.mimetype)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
