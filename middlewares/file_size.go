package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const max_upload_size = 1024 * 1024 //1MB

func BodySizeMiddleware(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	c.Request.Body = http.MaxBytesReader(w, c.Request.Body, max_upload_size)

	c.Next()
}
