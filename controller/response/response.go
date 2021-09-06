package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

func Success(c *gin.Context, data interface{}) {

	respond(c, http.StatusOK, response{Message: "SUCCESS", Data: data})
}

func Fail(c *gin.Context, err string) {
	respond(c, http.StatusAccepted, response{Message: "FAILED", Error: err})
}

func respond(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}
