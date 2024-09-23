package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func LoggingMiddleware(log *zap.SugaredLogger) gin.HandlerFunc {
	log.Infow("Zap logging started")
	return func(c *gin.Context) {
		start := time.Now()
		size := 0
		c.Next()
		end := time.Now().Sub(start) * time.Millisecond
		if c.Writer.Size() != -1 {
			size = c.Writer.Size()
		}

		log.Infow("request and response details",
			"url", c.Request.URL.String(),
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"size", size,
			"ms", end,
		)
	}
}
