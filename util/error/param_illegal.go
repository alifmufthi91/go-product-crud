package errorUtil

import (
	"net/http"
)

func ParamIllegal(message string) *CustomError {
	return &CustomError{
		Message:    message,
		HttpStatus: http.StatusBadRequest,
		ErrorName:  "PARAM_ILLEGAL",
	}
}
