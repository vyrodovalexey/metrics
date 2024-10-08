package sendmetrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/vyrodovalexey/metrics/internal/model"
	"io"
	"log"
	"net/http"
	"time"
)

// SendAsPlain Отправка запроса в формате plaintext
func SendAsPlain(cl *http.Client, url string) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err) // Ошибка при создании запроса
	}
	// Установка типа контента запроса
	req.Header.Set("Content-Type", "text/plain")

	resp, errr := cl.Do(req) // Отправка запроса

	if errr == nil {
		defer resp.Body.Close()

		// Вывод статуса запроса
		fmt.Println(time.Now(), " ", url, " ", resp.StatusCode)
	}
}

// SendAsJSON Отправка запроса в формате JSON
func SendAsJSON(cl *http.Client, url string, m *model.Metrics) {
	jm, _ := json.Marshal(*m)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jm))
	if err != nil {
		log.Println(err) // Ошибка при создании запроса
	}
	// Установка типа контента запроса и кодировок
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")
	// Пытаемся отправить запрос
	resp, errr := cl.Do(req)
	if errr == nil {
		defer resp.Body.Close()
		// Вывод статуса запроса
		fmt.Println(time.Now(), " ", url, " ", resp.StatusCode)
		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			// Обработка ответа, кодированного GZIP
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				log.Printf("failed to create gzip reader: %v", err) // Ошибка при создании читателя GZIP
			}
			defer reader.Close()
		default:
			// Ответ не кодирован GZIP, используем тело запроса как есть
			reader = resp.Body
		}
		body, err := io.ReadAll(reader) // Чтение тела ответа
		if err != nil {
			// Ошибка при чтении тела ответа
			fmt.Printf("Error reading response body: %v", err)
			return
		}

		// Вывод тела ответа
		log.Println(string(body))
	}
}
