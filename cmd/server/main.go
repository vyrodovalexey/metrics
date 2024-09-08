package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vyrodovalexey/metrics/internal/handlers"
	"github.com/vyrodovalexey/metrics/internal/storage"
)

func main() {

	gauge := make(map[string]storage.Gauge)
	counter := make(map[string][]storage.Counter)
	mst := storage.MemStorage{GaugeMap: gauge, CounterMap: counter}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.POST("/update/:type/:name/:value", handlers.Update(&mst))
	router.GET("/value/:type/:name", handlers.Get(&mst))
	router.GET("/", handlers.GetAllKeys(&mst))
	router.Run(":8080")

}
