package logging

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

func NewLogging(loglevel zapcore.Level) *zap.SugaredLogger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	loggerConfig.DisableCaller = true
	loggerConfig.Level.SetLevel(loglevel)

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	// nolint:errcheck
	defer logger.Sync()
	return logger.Sugar()
}

func LoggingMiddleware(log *zap.SugaredLogger) gin.HandlerFunc {
	log.Infow("Zap logging started")
	return func(c *gin.Context) {
		start := time.Now()
		size := 0
		c.Next()
		end := time.Since(start) * time.Millisecond
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
