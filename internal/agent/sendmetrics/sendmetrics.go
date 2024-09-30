package sendmetrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func SendAsPlain(cl *http.Client, url string) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, errr := cl.Do(req)

	if errr == nil {
		defer resp.Body.Close()
		fmt.Println(time.Now(), " ", url, " ", resp.StatusCode)
	}
}

func SendAsJSON(cl *http.Client, url string, m *Metrics) {
	jm, _ := json.Marshal(*m)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jm))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")
	// Attempt the request
	resp, errr := cl.Do(req)
	if errr == nil {
		defer resp.Body.Close()
		fmt.Println(time.Now(), " ", url, " ", resp.StatusCode)
		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			// Handle GZIP-encoded response
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				log.Printf("failed to create gzip reader: %w", err)
			}
			defer reader.Close()
		default:
			// Response is not gzipped, use the response body as is
			reader = resp.Body
		}
		body, err := io.ReadAll(reader)
		if err != nil {
			fmt.Printf("Error reading response body: %w", err)
			return
		}

		// Print the response body
		log.Println(string(body))
	}
}
