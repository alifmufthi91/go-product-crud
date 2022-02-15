package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Status  int         `json:"status,omitempty"`
}

func Success(c *gin.Context, data interface{}) {

	respond(c, http.StatusOK, response{Message: "SUCCESS", Data: data, Status: http.StatusOK})
}

func Fail(c *gin.Context, err string) {
	respond(c, http.StatusInternalServerError, response{Message: "Internal Error", Error: err, Status: http.StatusInternalServerError})
}

func respond(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}
