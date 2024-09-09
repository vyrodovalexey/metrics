package main

import (
	"github.com/vyrodovalexey/metrics/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	gauge := make(map[string]storage.Gauge)
	counter := make(map[string]storage.Counter)
	mst := storage.MemStorage{GaugeMap: gauge, CounterMap: counter}
	router := SetupRouter(&mst)
	router.LoadHTMLGlob("../../templates/*")
	tests := []struct {
		name           string
		st             *storage.MemStorage
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
		/*		{
				name:           "Invalid MediaType",
				st:             &mst,
				method:         http.MethodPost,
				url:            "/update/gauge/test/1.12",
				mimetype:       "application/json",
				expectedStatus: http.StatusUnsupportedMediaType,
			},*/
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
