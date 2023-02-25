package responseUtil

import (
	"net/http"
	"product-crud/dto/app"
	errorUtil "product-crud/util/error"
	"product-crud/util/logger"

	"github.com/gin-gonic/gin"
)

func ErrorHandling(c *gin.Context) {
	if r := recover(); r != nil {
		logger.Error("Recovered from panic: %+v", r)
		var errorMessage string
		if err, ok := r.(*errorUtil.CustomError); ok {
			Fail(c, app.ErrorHttpResponse{
				Message:    err.Error(),
				HttpStatus: err.HttpStatus,
				ErrorName:  err.ErrorName,
			})
			return
		} else if err, ok := r.(error); ok {
			errorMessage = err.Error()
		} else if err, ok := r.(string); ok {
			errorMessage = err
		} else {
			errorMessage = "Internal Error"
		}
		Fail(c, app.ErrorHttpResponse{
			Message:    errorMessage,
			HttpStatus: http.StatusInternalServerError,
			ErrorName:  "INTERNAL SERVER ERROR",
		})
		return
	}
}
