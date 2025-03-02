package hash

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func VerifyHashMiddleware(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		receivedHash := c.GetHeader("HashSHA256")
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		hash := sha256.New()
		hash.Write(bodyBytes)
		hash.Write([]byte(key))
		computedHash := hex.EncodeToString(hash.Sum(nil))

		if computedHash != receivedHash && receivedHash != "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Next()
	}
}

func AddResponseHashMiddleware(key string) gin.HandlerFunc {
	return func(c *gin.Context) {

		writer := &bodyWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		body := writer.body.Bytes()
		hash := sha256.New()
		hash.Write(body)
		hash.Write([]byte(key))
		computedHash := hex.EncodeToString(hash.Sum(nil))
		c.Writer.Header().Set("HashSHA256", computedHash)
	}
}
