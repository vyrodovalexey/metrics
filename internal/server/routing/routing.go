package routing

import (
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/vyrodovalexey/metrics/internal/server/handlers"
	"github.com/vyrodovalexey/metrics/internal/server/logging"
	"github.com/vyrodovalexey/metrics/internal/server/storage"
	"go.uber.org/zap"
	"io"
)

func SetupRouter(log *zap.SugaredLogger) *gin.Engine {
	// Установка режима работы Gin в release-режиме
	gin.SetMode(gin.ReleaseMode)
	// Установка стандартного вывода Gin в discard-режиме
	gin.DefaultWriter = io.Discard
	router := gin.Default()
	// Добавление middleware для логирования
	router.Use(logging.LoggingMiddleware(log))
	// Добавление middleware для сжатия ответа
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	return router
}

func ConfigureRouting(r *gin.Engine, st storage.Storage) {
	// Определение эндпоинтов
	//r.GET("/ping", handlers.CheckDatabaseConnection(st))
	r.POST("/update/:type/:name/:value", handlers.UpdateFromURLPath(st))
	r.GET("/value/:type/:name", handlers.Get(st))
	r.POST("/update/", handlers.UpdateFromBodyJSON(st))
	r.POST("/value/", handlers.GetBodyJSON(st))
	r.GET("/", handlers.GetAllKeys(st))
}
