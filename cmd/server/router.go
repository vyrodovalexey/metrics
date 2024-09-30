package main

import (
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/vyrodovalexey/metrics/internal/handlers"
	"github.com/vyrodovalexey/metrics/internal/server/logging"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"go.uber.org/zap"
	"io"
)

func SetupRouter(st storage.Storage, log *zap.SugaredLogger) *gin.Engine {
	// Установка режима работы Gin в release-режиме
	gin.SetMode(gin.ReleaseMode)
	// Установка стандартного вывода Gin в discard-режиме
	gin.DefaultWriter = io.Discard
	router := gin.Default()
	// Добавление middleware для логирования
	router.Use(logging.LoggingMiddleware(log))
	// Добавление middleware для сжатия ответа
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	// Определение эндпоинтов
	router.POST("/update/:type/:name/:value", handlers.Update(st))
	router.GET("/value/:type/:name", handlers.Get(st))
	router.POST("/update/", handlers.UpdateJSON(st))
	router.POST("/value/", handlers.GetJSON(st))
	router.GET("/", handlers.GetAllKeys(st))
	return router
}
