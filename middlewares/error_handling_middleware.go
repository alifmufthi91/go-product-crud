package middlewares

import (
	"net/http"
	"product-crud/dto/app"
	"product-crud/util/apiresponse"
	"product-crud/util/errorhandler"
	"product-crud/util/logger"

	"github.com/gin-gonic/gin"
)

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Recovered from panic: %+v", r)
				var errorMessage string
				if err, ok := r.(*errorhandler.CustomError); ok {
					apiresponse.Fail(c, app.ErrorHttpResponse{
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
				apiresponse.Fail(c, app.ErrorHttpResponse{
					Message:    errorMessage,
					HttpStatus: http.StatusInternalServerError,
					ErrorName:  "INTERNAL SERVER ERROR",
				})
			}
		}()
		c.Next()
	}
}
