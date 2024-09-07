package handlers

import (
	"github.com/vyrodovalexey/metrics/internal/storage"
	"net/http"
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
		if len(pathSlice) != 4 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch pathSlice[1] {
		case "gauge":
			err := st.AddGauge(pathSlice[2], pathSlice[3])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case "counter":
			err := st.AddCounter(pathSlice[2], pathSlice[3])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
