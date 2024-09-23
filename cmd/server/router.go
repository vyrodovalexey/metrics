package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vyrodovalexey/metrics/internal/handlers"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"go.uber.org/zap"
	"io"
)

func SetupRouter(mst *storage.MemStorage, log *zap.SugaredLogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	router := gin.Default()
	router.Use(LoggingMiddleware(log))
	router.POST("/update/:type/:name/:value", handlers.Update(mst))
	router.GET("/value/:type/:name", handlers.Get(mst))
	router.GET("/", handlers.GetAllKeys(mst))
	return router
}
