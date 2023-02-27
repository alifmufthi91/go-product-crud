package util

import (
	"product-crud/dto/app"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GeneratePaginationFromRequest(c *gin.Context) (*app.Pagination, error) {
	// Initializing default
	//	var mode string
	limit := 5
	page := 1
	sort := "created_at asc"
	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			val, err := strconv.Atoi(queryValue)
			if err != nil {
				return nil, err
			}
			limit = val
		case "page":
			val, err := strconv.Atoi(queryValue)
			if err != nil {
				return nil, err
			}
			page = val
		case "sort":
			sort = queryValue
		}
	}
	return &app.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}, nil

}
