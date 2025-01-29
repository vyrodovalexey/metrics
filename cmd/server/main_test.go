package main

import (
	"bytes"
	"context"
	"github.com/vyrodovalexey/metrics/internal/server/config"
	"github.com/vyrodovalexey/metrics/internal/server/logging"
	"github.com/vyrodovalexey/metrics/internal/server/memstorage"
	"github.com/vyrodovalexey/metrics/internal/server/routing"
	storage2 "github.com/vyrodovalexey/metrics/internal/server/storage"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg := config.New()
	ConfigParser(cfg)
}

func TestRequestsMemStorageSyncNew(t *testing.T) {
	var st storage2.Storage = &memstorage.MemStorageWithAttributes{}
	ctx := context.Background()
	sugar := logging.NewLogging(zap.InfoLevel)
	// Создаем новый экземпляр конфигурации
	err := st.New(ctx, "/tmp/metrics-storage.json", 0)
	if err != nil {
		t.Errorf("initializing file storage... Error: %v", err)
	}

	router := routing.SetupRouter(sugar)
	routing.ConfigureRouting(ctx, router, st)
	router.LoadHTMLGlob("../../templates/*")

	testsNew := []struct {
		name           string
		method         string
		url            string
		mimetype       string
		body           string
		expectedStatus int
		expectedValue  string
	}{
		{
			name:           "Valid Update gauge",
			method:         http.MethodPost,
			url:            "/update/gauge/test/1.4545",
			mimetype:       "text/plain",
			expectedStatus: http.StatusOK,
			expectedValue:  "1.4545",
		},
		{
			name:           "Valid Update counter",
			method:         http.MethodPost,
			url:            "/update/counter/test/1",
			mimetype:       "text/plain",
			expectedStatus: http.StatusOK,
			expectedValue:  "1",
		},
		{
			name:           "Invalid Method",
			method:         http.MethodGet,
			url:            "/update/gauge/test/1.12",
			mimetype:       "text/plain",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid Gauge",
			method:         http.MethodPost,
			url:            "/update/gauge/test/test",
			mimetype:       "text/plain",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Counter",
			method:         http.MethodPost,
			url:            "/update/counter/test/1.12",
			mimetype:       "text/plain",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Get counter",
			method:         http.MethodGet,
			url:            "/value/counter/test",
			mimetype:       "text/plain",
			expectedStatus: http.StatusOK,
			expectedValue:  "1",
		},
		{
			name:           "Get gauge",
			method:         http.MethodGet,
			url:            "/value/gauge/test",
			mimetype:       "text/plain",
			expectedStatus: http.StatusOK,
			expectedValue:  "1.4545",
		},
		{
			name:           "Invalid Get gauge",
			method:         http.MethodGet,
			url:            "/value/gauge/unavailable",
			mimetype:       "text/plain",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Get /",
			method:         http.MethodGet,
			url:            "/",
			mimetype:       "text/html",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Post Counter Json /update",
			method:         http.MethodPost,
			url:            "/update/",
			mimetype:       "application/json",
			body:           "{\"id\":\"test\",\"type\":\"counter\",\"delta\":1}",
			expectedStatus: http.StatusOK,
			expectedValue:  "{\"id\":\"test\",\"type\":\"counter\",\"delta\":2}",
		},
		{
			name:           "Post Counter Json /value",
			method:         http.MethodPost,
			url:            "/value/",
			mimetype:       "application/json",
			body:           "{\"id\":\"test\",\"type\":\"counter\",\"delta\":5}",
			expectedStatus: http.StatusOK,
			expectedValue:  "{\"id\":\"test\",\"type\":\"counter\",\"delta\":2}",
		},
		{
			name:           "Post Gauge Json /update",
			method:         http.MethodPost,
			url:            "/update/",
			mimetype:       "application/json",
			body:           "{\"id\":\"testjson\",\"type\":\"gauge\",\"value\":1.678}",
			expectedStatus: http.StatusOK,
			expectedValue:  "{\"id\":\"testjson\",\"type\":\"gauge\",\"value\":1.678}",
		},
		{
			name:           "Post Gauge Json /value",
			method:         http.MethodPost,
			url:            "/value/",
			mimetype:       "application/json",
			body:           "{\"id\":\"testjson\",\"type\":\"gauge\"}",
			expectedStatus: http.StatusOK,
			expectedValue:  "{\"id\":\"testjson\",\"type\":\"gauge\",\"value\":1.678}",
		},
		{
			name:           "Post Json Batch /updates/",
			method:         http.MethodPost,
			url:            "/updates/",
			mimetype:       "application/json",
			body:           "[{\"id\":\"test\",\"type\":\"counter\",\"delta\":1},{\"id\":\"testbatch\",\"type\":\"gauge\",\"value\":1.5}]",
			expectedStatus: http.StatusOK,
			expectedValue:  "{\"id\":\"test\",\"type\":\"counter\",\"delta\":3}{\"id\":\"testbatch\",\"type\":\"gauge\",\"value\":1.5}",
		},
		{
			name:           "Post Gauge Json /value",
			method:         http.MethodPost,
			url:            "/value/",
			mimetype:       "application/json",
			body:           "{\"id\":\"testbatch\",\"type\":\"gauge\"}",
			expectedStatus: http.StatusOK,
			expectedValue:  "{\"id\":\"testbatch\",\"type\":\"gauge\",\"value\":1.5}",
		},
	}

	for _, tt := range testsNew {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.body == "" {
				body = nil
			} else {
				body = bytes.NewBuffer([]byte(tt.body))
			}
			req := httptest.NewRequest(tt.method, tt.url, body)
			req.Header.Add("Content-Type", tt.mimetype)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if tt.expectedValue != "" {
				if w.Body.String() != tt.expectedValue {
					t.Errorf("expected value %s, got %s", tt.expectedValue, w.Body.String())
				} else {
					t.Logf("expected value %s, got %s", tt.expectedValue, w.Body.String())
				}
			}
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

		})
	}
	st.Close()

}

func TestMemStorageSyncLoad(t *testing.T) {
	var st storage2.Storage = &memstorage.MemStorageWithAttributes{}
	ctx := context.Background()
	sugar := logging.NewLogging(zap.InfoLevel)
	err := st.Load(ctx, "../../test/data/metrics-storage.json", 0)
	if err != nil {
		// Логируем ошибку, если открытие/создание файла не удалось
		sugar.Panicw("Initializing file storage...",
			"Error load data from file:", err,
		)
		return
	}

	router := routing.SetupRouter(sugar)
	routing.ConfigureRouting(ctx, router, st)
	router.LoadHTMLGlob("../../templates/*")

	testsLoad := []struct {
		name           string
		method         string
		url            string
		mimetype       string
		body           string
		expectedStatus int
		expectedValue  string
	}{
		{
			name:           "Post Counter Json /value",
			method:         http.MethodPost,
			url:            "/value/",
			mimetype:       "application/json",
			body:           "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":5}",
			expectedStatus: http.StatusOK,
			expectedValue:  "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":1}",
		},
		{
			name:           "Post Counter Json /update",
			method:         http.MethodPost,
			url:            "/update/",
			mimetype:       "application/json",
			body:           "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":5}",
			expectedStatus: http.StatusOK,
			expectedValue:  "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":6}",
		},
		{
			name:           "Post Gauge Json /value",
			method:         http.MethodPost,
			url:            "/value/",
			mimetype:       "application/json",
			body:           "{\"id\":\"TestGauge\",\"type\":\"gauge\"}",
			expectedStatus: http.StatusOK,
			expectedValue:  "{\"id\":\"TestGauge\",\"type\":\"gauge\",\"value\":2.4545}",
		},
		{
			name:           "Post Gauge Json /update",
			method:         http.MethodPost,
			url:            "/update/",
			mimetype:       "application/json",
			body:           "{\"id\":\"TestGauge\",\"type\":\"gauge\",\"value\":1.678}",
			expectedStatus: http.StatusOK,
			expectedValue:  "{\"id\":\"TestGauge\",\"type\":\"gauge\",\"value\":1.678}",
		},
	}

	for _, tt := range testsLoad {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.body == "" {
				body = nil
			} else {
				body = bytes.NewBuffer([]byte(tt.body))
			}
			req := httptest.NewRequest(tt.method, tt.url, body)
			req.Header.Add("Content-Type", tt.mimetype)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if tt.expectedValue != "" {
				if w.Body.String() != tt.expectedValue {
					t.Errorf("expected value %s, got %s", tt.expectedValue, w.Body.String())
				} else {
					t.Logf("expected value %s, got %s", tt.expectedValue, w.Body.String())
				}
			}
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

		})
	}
	st.Close()
}

// Temp remove pgstorage tests
//func TestPgStorage(t *testing.T) {
//	ctx := context.Background()
//	pg, _ := sqltestutil.StartPostgresContainer(ctx, "14")
//	defer pg.Shutdown(ctx)
//
//	var st storage2.Storage = &pgstorage.PgStorageWithAttributes{}
//	fmt.Println(pg.ConnectionString())
//	sugar := logging.NewLogging(zap.InfoLevel)
//	err := st.New(ctx, pg.ConnectionString(), 0)
//
//	if err != nil {
//		// Логируем ошибку, если открытие/создание файла не удалось
//		sugar.Panicw("Connecting to database...",
//			"Error database connection:", err,
//		)
//		return
//	}
//	sugar.Infow("Connected to database")
//
//	router := routing.SetupRouter(sugar)
//	routing.ConfigureRouting(ctx, router, st)
//	router.LoadHTMLGlob("../../templates/*")
//
//	testsLoad := []struct {
//		name           string
//		method         string
//		url            string
//		mimetype       string
//		body           string
//		expectedStatus int
//		expectedValue  string
//	}{
//		{
//			name:           "Post Counter Json /update",
//			method:         http.MethodPost,
//			url:            "/update/",
//			mimetype:       "application/json",
//			body:           "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":7}",
//			expectedStatus: http.StatusOK,
//			expectedValue:  "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":7}",
//		},
//		{
//			name:           "Post Gauge Json /update",
//			method:         http.MethodPost,
//			url:            "/update/",
//			mimetype:       "application/json",
//			body:           "{\"id\":\"TestGauge\",\"type\":\"gauge\",\"value\":1.678}",
//			expectedStatus: http.StatusOK,
//			expectedValue:  "{\"id\":\"TestGauge\",\"type\":\"gauge\",\"value\":1.678}",
//		},
//		{
//			name:           "Post Counter Json /value",
//			method:         http.MethodPost,
//			url:            "/value/",
//			mimetype:       "application/json",
//			body:           "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":8}",
//			expectedStatus: http.StatusOK,
//			expectedValue:  "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":7}",
//		},
//		{
//			name:           "Post Gauge Json /value",
//			method:         http.MethodPost,
//			url:            "/value/",
//			mimetype:       "application/json",
//			body:           "{\"id\":\"TestGauge\",\"type\":\"gauge\"}",
//			expectedStatus: http.StatusOK,
//			expectedValue:  "{\"id\":\"TestGauge\",\"type\":\"gauge\",\"value\":1.678}",
//		},
//		{
//			name:           "Get /",
//			method:         http.MethodGet,
//			url:            "/",
//			mimetype:       "text/html",
//			expectedStatus: http.StatusOK,
//		},
//		{
//			name:           "Get /ping",
//			method:         http.MethodGet,
//			url:            "/ping",
//			mimetype:       "text/html",
//			expectedStatus: http.StatusOK,
//		},
//	}
//
//	for _, tt := range testsLoad {
//		t.Run(tt.name, func(t *testing.T) {
//			var body io.Reader
//			if tt.body == "" {
//				body = nil
//			} else {
//				body = bytes.NewBuffer([]byte(tt.body))
//			}
//			req := httptest.NewRequest(tt.method, tt.url, body)
//			req.Header.Add("Content-Type", tt.mimetype)
//			w := httptest.NewRecorder()
//			router.ServeHTTP(w, req)
//			if tt.expectedValue != "" {
//				if w.Body.String() != tt.expectedValue {
//					t.Errorf("expected value %s, got %s", tt.expectedValue, w.Body.String())
//				} else {
//					t.Logf("expected value %s, got %s", tt.expectedValue, w.Body.String())
//				}
//			}
//			if w.Code != tt.expectedStatus {
//				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
//			}
//
//		})
//	}
//	st.Close()
//}
