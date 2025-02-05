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

const (
	ContentType          = "Content-Type"
	ContentEncoding      = "Content-Encoding"
	ContentTypeTextPlain = "text/plain"
	ContentTypeJSON      = "application/json"
	EncodingGzip         = "gzip"
)

func SendRequest(cl *http.Client, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i < 3; i++ {
		resp, err = cl.Do(req)

		if err != nil {
			fmt.Printf("Server is not ready: %v\n", err)
		} else {
			return resp, nil
		}
		if i == 0 {
			<-time.After(1 * time.Second)
		} else {
			<-time.After(time.Duration(i*2+1) * time.Second)
		}
	}

	return resp, err
}

func SendAsJSONRequest(cl *http.Client, url string, jm []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jm))
	if err != nil {
		log.Println(err) // Ошибка при создании запроса
	}
	// Установка типа контента запроса и кодировок
	req.Header.Set(ContentType, ContentTypeJSON)
	req.Header.Set("Accept-Encoding", EncodingGzip)
	req.Header.Set(ContentEncoding, EncodingGzip)
	// Пытаемся отправить запрос
	resp, errr := SendRequest(cl, req)
	if errr == nil {

		// Вывод статуса запроса
		fmt.Println(time.Now(), " ", url, " ", resp.StatusCode)
		var reader io.ReadCloser
		switch resp.Header.Get(ContentEncoding) {
		case EncodingGzip:
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
		resp.Body.Close()
		body, err := io.ReadAll(reader) // Чтение тела ответа
		if err != nil {
			// Ошибка при чтении тела ответа
			fmt.Printf("Error reading response body: %v", err)
			return err
		}

		// Вывод тела ответа
		log.Println(string(body))

	}
	return errr
}

// SendAsPlain Отправка запроса в формате plaintext
func SendAsPlain(cl *http.Client, url string) error {
	//timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	//defer cancel()
	req, err := http.NewRequest("POST", url, http.NoBody)
	if err != nil {
		log.Printf("failed to create request: %v", err) // Ошибка при создании запроса
	}
	// Установка типа контента запроса
	req.Header.Set(ContentType, ContentTypeTextPlain)

	resp, errr := SendRequest(cl, req) // Отправка запроса

	if errr == nil {

		// Вывод статуса запроса
		fmt.Println(time.Now(), " ", url, " ", resp.StatusCode)
		resp.Body.Close()
	}
	return errr
}

// SendAsJSON Отправка запроса в формате JSON
func SendAsJSON(cl *http.Client, url string, m *model.Metrics) error {
	jm, _ := json.Marshal(*m)
	//timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	//defer cancel()
	err := SendAsJSONRequest(cl, url, jm)
	return err

}

// SendAsBatchJSON Отправка batch в формате JSON
func SendAsBatchJSON(cl *http.Client, url string, b *model.MetricsBatch) error {
	jm, _ := json.Marshal(*b)
	//timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	//defer cancel()
	err := SendAsJSONRequest(cl, url, jm)
	return err
}
