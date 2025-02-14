package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vyrodovalexey/metrics/internal/model"
	"github.com/vyrodovalexey/metrics/internal/server/storage"
	"net/http"
)

const (
	badrequest      = "Bad Request"
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
	ContentEncoding = "Content-Encoding"
	EncodingGzip    = "gzip"
)

// Attributes структура для хранения информации о хранилище
type Attributes struct {
	St storage.Storage
}

// UpdateFromBodyJSON обновляет метрику из тела запроса в формате JSON.
func (a *Attributes) UpdateFromBodyJSON(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, что Content-Type запроса - application/json
		if c.Request.Header.Get(ContentType) != ContentTypeJSON {
			// Если нет, возвращаем ошибку 415 Unsupported Media Type
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error": badrequest,
			})
			return
		} else {
			// Создаем новую пустую метрику
			m := &model.Metrics{}
			// Получаем тело запроса
			body := c.Request.Body
			// Парсим тело запроса в структуру Metrics
			err := m.BodyToMetric(body)
			// Устанавливаем заголовок Content-Type в application/json
			c.Header(ContentType, ContentTypeJSON)
			c.Header(ContentEncoding, EncodingGzip)
			// Если произошла ошибка при парсинге, возвращаем ошибку 400 Bad Request
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
				return
			}
			// Обновляем метрику в хранилище
			err = a.St.UpdateMetric(ctx, m)
			// Если произошла ошибка при обновлении, возвращаем ошибку 500 Internal Server Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			// Получаем обновленную метрику из хранилища
			a.St.GetMetric(ctx, m)
			// Возвращаем обновленную метрику клиенту с кодом 200 OK
			c.JSON(http.StatusOK, m)
			return
		}
	}
}

// UpdateFromBodyJSON обновляет метрику из тела запроса в формате JSON.
func (a *Attributes) BatchUpdateFromBodyJSON(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, что Content-Type запроса - application/json
		if c.Request.Header.Get(ContentType) != ContentTypeJSON {
			// Если нет, возвращаем ошибку 415 Unsupported Media Type
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error": badrequest,
			})
			return
		} else {
			// Создаем новую пустую метрику
			b := &model.MetricsBatch{}
			// Получаем тело запроса
			body := c.Request.Body
			// Парсим тело запроса в структуру Metrics
			err := b.BodyToMetricBatch(body)
			// Устанавливаем заголовок Content-Type в application/json
			c.Header(ContentType, ContentTypeJSON)
			c.Header(ContentEncoding, EncodingGzip)
			// Если произошла ошибка при парсинге, возвращаем ошибку 400 Bad Request
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
				return
			}
			for i := range *b {
				// Обновляем метрику в хранилище
				err = a.St.UpdateMetric(ctx, &(*b)[i])
				// Если произошла ошибка при обновлении, возвращаем ошибку 500 Internal Server Error
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}
				// Получаем обновленную метрику из хранилища
				a.St.GetMetric(ctx, &(*b)[i])
				// Возвращаем обновленную метрику клиенту с кодом 200 OK
				c.JSON(http.StatusOK, &(*b)[i])
			}
			return
		}
	}
}

// UpdateFromURLPath обновляет метрику из параметров URL.
func (a *Attributes) UpdateFromURLPath(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Создаем новую пустую метрику
		m := &model.Metrics{}
		// Парсим параметры URL в структуру Metrics
		err := m.URLPathToMetric(c.Param("type"), c.Param("name"), c.Param("value"))
		// Если произошла ошибка при парсинге, возвращаем ошибку 400
		if err != nil {
			c.String(http.StatusBadRequest, badrequest)
			return
		}
		// Обновляем метрику в хранилище
		err = a.St.UpdateMetric(ctx, m)
		// Если произошла ошибка при обновлении, возвращаем ошибку 400
		if err != nil {
			c.String(http.StatusBadRequest, badrequest)
			return
		}
		// Получаем обновленную метрику из хранилища
		a.St.GetMetric(ctx, m)
		// Возвращаем обновленную метрику клиенту с кодом 200 OK
		c.String(http.StatusOK, m.String())
	}
}

// CheckDatabaseConnection проверяет соединение с базой данных.
func (a *Attributes) CheckDatabaseConnection(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := a.St.Check(ctx)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
		} else {
			c.String(http.StatusOK, "ok")
		}
	}
}

// Get возвращает метрику по ее имени.
func (a *Attributes) Get(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Создаем новую пустую метрику
		m := &model.Metrics{}
		// Парсим параметры URL в структуру Metrics
		err := m.URLPathToMetric(c.Param("type"), c.Param("name"), c.Param("value"))
		// Если произошла ошибка при парсинге, возвращаем ошибку 400
		if err != nil {
			c.String(http.StatusBadRequest, badrequest)
			return
		}
		// Получаем метрику из хранилища
		b := a.St.GetMetric(ctx, m)
		// Если метрика не найдена, возвращаем ошибку 404 Not Found
		if !b {
			c.String(http.StatusNotFound, badrequest)
			return
		}
		// Возвращаем метрику клиенту с кодом 200 OK
		c.String(http.StatusOK, m.String())
	}
}

// GetBodyJSON возвращает метрику из тела запроса в формате JSON.
func (a *Attributes) GetBodyJSON(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, что Content-Type запроса - application/json
		if c.Request.Header.Get(ContentType) != ContentTypeJSON {
			// Если нет, возвращаем ошибку 415 Unsupported Media Type
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error": badrequest,
			})
			return
		}
		// Создаем новую пустую метрику
		m := &model.Metrics{}
		c.Header(ContentType, ContentTypeJSON)
		c.Header(ContentEncoding, EncodingGzip)
		// Получаем тело запроса
		body := c.Request.Body
		// Парсим тело запроса в структуру Metrics
		err := m.BodyToMetric(body)
		// Если произошла ошибка при парсинге, возвращаем ошибку 400 Bad Request
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": badrequest,
			})
			return
		}
		// Получаем метрику из хранилища
		b := a.St.GetMetric(ctx, m)
		// Если метрика не найдена, возвращаем ошибку 404 Not Found
		if !b {
			c.JSON(http.StatusNotFound, gin.H{
				"error": badrequest,
			})
			return
		}
		// Возвращаем метрику клиенту с кодом 200 OK
		c.JSON(http.StatusOK, m)
	}
}

// GetAllKeys возвращает все ключи метрик.
func (a *Attributes) GetAllKeys(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем все ключи метрик из хранилища
		gval, cval, err := a.St.GetAllMetricNames(ctx)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
		}
		// Возвращаем ключи метрик клиенту с кодом 200 OK
		c.HTML(http.StatusOK, "table.tmpl", gin.H{
			"Title":         "Metric Names",
			"GaugeValues":   gval,
			"CounterValues": cval,
		})

	}
}
