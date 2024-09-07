package handlers

import (
	"github.com/vyrodovalexey/metrics/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

type StorageHandler struct {
	handler http.Handler
	store   storage.MemStorage
}

// ServeHTTP handles the request by passing it to the real
// handler and logging the request details
func (sh *StorageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sh.handler.ServeHTTP(w, r)

	pathSlice := strings.Split(r.URL.Path[1:], "/")

	if len(pathSlice) == 3 && (pathSlice[1] == "gauge" || pathSlice[1] == "counter") {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(pathSlice) != 4 || (pathSlice[1] != "gauge" && pathSlice[1] != "counter") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if pathSlice[1] == "gauge" {
		_, err := strconv.ParseFloat(pathSlice[3], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			gauge, _ := strconv.ParseFloat(pathSlice[3], 64)
			sh.store.AddGauge(pathSlice[2], gauge)
		}
	}

	if pathSlice[1] == "counter" {
		_, err := strconv.ParseInt(pathSlice[3], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			counter, _ := strconv.ParseInt(pathSlice[3], 10, 64)
			sh.store.AddCounter(pathSlice[2], counter)
		}
	}
}

// NewLogger constructs a new Logger middleware handler
func NewStorageHandler(handlerToWrap http.Handler, st *storage.MemStorage) *StorageHandler {
	return &StorageHandler{handlerToWrap, *st}
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

}
