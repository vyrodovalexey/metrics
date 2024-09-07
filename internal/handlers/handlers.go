package handlers

import (
	"github.com/vyrodovalexey/metrics/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

func Update(st *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "text/plain" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

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
				st.AddGauge(pathSlice[2], gauge)
			}
		}

		if pathSlice[1] == "counter" {
			_, err := strconv.ParseInt(pathSlice[3], 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				counter, _ := strconv.ParseInt(pathSlice[3], 10, 64)
				st.AddCounter(pathSlice[2], counter)
			}
		}
	}
}
