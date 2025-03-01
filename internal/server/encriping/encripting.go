package encriping

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	badrequest = "Bad Request"
)

func CheckShaSumHeader(shasum [32]byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("HashSHA256") != fmt.Sprintf("%x", shasum) {
			c.String(http.StatusBadRequest, badrequest)
			c.Abort()
			return
		}
		c.Next()
		c.Header("HashSHA256", fmt.Sprintf("%x", shasum))
	}
}
