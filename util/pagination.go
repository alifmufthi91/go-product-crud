package util

import (
	"product-crud/app"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GeneratePaginationFromRequest(c *gin.Context) app.Pagination {
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
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		case "sort":
			sort = queryValue
		}
	}
	return app.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}

}
