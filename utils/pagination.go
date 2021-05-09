package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mop/entity"
)

//GeneratePaginationFromRequest ..
func GeneratePaginationFromRequest(c *gin.Context) entity.Pagination {
	// Initializing default
	//	var mode string
	limit := 20
	page := 1
	sort := "created_at desc"
	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "sort":
			sort = queryValue
			break

		}
	}
	return entity.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}

}
