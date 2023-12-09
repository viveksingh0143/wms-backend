package middlewares

import (
	"github.com/gin-gonic/gin"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pagination requests.Pagination
		var sorting requests.Sorting

		if err := c.ShouldBindQuery(&pagination); err != nil {
			resp := responses.NewValidationErrorResponse(err, pagination)
			c.JSON(resp.Status, resp)
			return
		}
		pagination.CalculatePageAndPageSize()

		if err := c.ShouldBindQuery(&sorting); err != nil {
			resp := responses.NewValidationErrorResponse(err, sorting)
			c.JSON(resp.Status, resp)
			return
		}
		sorting.CalculateSorting()

		c.Set("pagination", pagination)
		c.Set("sorting", sorting)
		c.Next()
	}
}
