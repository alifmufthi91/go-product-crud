package errorUtil

import (
	"net/http"
)

func DataNotFound(message string) *CustomError {
	return &CustomError{
		Message:    message,
		HttpStatus: http.StatusNotFound,
		ErrorName:  "NOT FOUND",
	}
}
