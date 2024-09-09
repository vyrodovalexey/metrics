package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vyrodovalexey/metrics/internal/handlers"
	"github.com/vyrodovalexey/metrics/internal/storage"
)

func SetupRouter(mst *storage.MemStorage) *gin.Engine {
	router := gin.Default()
	router.POST("/update/:type/:name/:value", handlers.Update(mst))
	router.GET("/value/:type/:name", handlers.Get(mst))
	router.GET("/", handlers.GetAllKeys(mst))
	return router
}
