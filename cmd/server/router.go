package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vyrodovalexey/metrics/internal/handlers"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"go.uber.org/zap"
	"io"
)

func SetupRouter(st storage.Storage, log *zap.SugaredLogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	router := gin.Default()
	router.Use(LoggingMiddleware(log))
	router.POST("/update/:type/:name/:value", handlers.Update(st))
	router.GET("/value/:type/:name", handlers.Get(st))
	router.POST("/update/", handlers.UpdateJson(st))
	router.POST("/value/", handlers.GetJson(st))
	router.GET("/", handlers.GetAllKeys(st))
	return router
}
