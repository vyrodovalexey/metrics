package hash

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func CheckShaSumHeader(key string) gin.HandlerFunc {
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
