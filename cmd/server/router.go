package main

import (
	"context"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/vyrodovalexey/metrics/internal/handlers"
	"github.com/vyrodovalexey/metrics/internal/server/logging"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"go.uber.org/zap"
	"io"
	"os"
)

func SetupRouter(ctx context.Context, st storage.Storage, conn *pgx.Conn, f *os.File, log *zap.SugaredLogger, p bool) *gin.Engine {
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
	router.GET("/ping", handlers.CheckDatabaseConnection(ctx, conn))
	router.POST("/update/:type/:name/:value", handlers.UpdateFromURLPath(st, f, p))
	router.GET("/value/:type/:name", handlers.Get(st))
	router.POST("/update/", handlers.UpdateFromBodyJSON(st, f, p))
	router.POST("/value/", handlers.GetBodyJSON(st))
	router.GET("/", handlers.GetAllKeys(st))
	return router
}
