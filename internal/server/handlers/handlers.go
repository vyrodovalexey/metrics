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
	badrequest = "Bad Request"
)

// UpdateFromBodyJSON обновляет метрику из тела запроса в формате JSON.
func UpdateFromBodyJSON(ctx context.Context, st storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, что Content-Type запроса - application/json
		if c.Request.Header.Get("Content-Type") != "application/json" {
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
			err := m.BodyToMetric(&body)
			// Если произошла ошибка при парсинге, возвращаем ошибку 400 Bad Request
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
				return
			}
			// Обновляем метрику в хранилище
			err = st.UpdateMetric(ctx, m)
			// Если произошла ошибка при обновлении, возвращаем ошибку 500 Internal Server Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			// Получаем обновленную метрику из хранилища
			st.GetMetric(ctx, m)
			// Возвращаем обновленную метрику клиенту с кодом 200 OK
			c.JSON(http.StatusOK, m)
			return
		}
	}
}

// UpdateFromURLPath обновляет метрику из параметров URL.
func UpdateFromURLPath(ctx context.Context, st storage.Storage) gin.HandlerFunc {
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
		err = st.UpdateMetric(ctx, m)
		// Если произошла ошибка при обновлении, возвращаем ошибку 400
		if err != nil {
			c.String(http.StatusBadRequest, badrequest)
			return
		}
		// Получаем обновленную метрику из хранилища
		st.GetMetric(ctx, m)
		// Возвращаем обновленную метрику клиенту с кодом 200 OK
		c.String(http.StatusOK, m.PrintMetric())
	}
}

// CheckDatabaseConnection проверяет соединение с базой данных.
func CheckDatabaseConnection(ctx context.Context, st storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := st.Check(ctx)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
		}
		c.String(http.StatusOK, "ok")
	}
}

// Get возвращает метрику по ее имени.
func Get(ctx context.Context, st storage.Storage) gin.HandlerFunc {
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
		b := st.GetMetric(ctx, m)
		// Если метрика не найдена, возвращаем ошибку 404 Not Found
		if !b {
			c.String(http.StatusNotFound, badrequest)
			return
		}
		// Возвращаем метрику клиенту с кодом 200 OK
		c.String(http.StatusOK, m.PrintMetric())

	}
}

// GetBodyJSON возвращает метрику из тела запроса в формате JSON.
func GetBodyJSON(ctx context.Context, st storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, что Content-Type запроса - application/json
		if c.Request.Header.Get("Content-Type") != "application/json" {
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
			err := m.BodyToMetric(&body)
			// Если произошла ошибка при парсинге, возвращаем ошибку 400 Bad Request
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
				return
			}
			// Получаем метрику из хранилища
			b := st.GetMetric(ctx, m)
			// Если метрика не найдена, возвращаем ошибку 404 Not Found
			if !b {
				c.JSON(http.StatusNotFound, gin.H{
					"error": badrequest,
				})
				return
			}
			// Возвращаем метрику клиенту с кодом 200 OK
			c.JSON(http.StatusOK, m)
			return
		}
	}
}

// GetAllKeys возвращает все ключи метрик.
func GetAllKeys(ctx context.Context, st storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем все ключи метрик из хранилища
		gval, cval, err := st.GetAllMetricNames(ctx)
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
