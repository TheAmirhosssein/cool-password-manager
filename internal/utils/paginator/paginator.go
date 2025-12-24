package paginator

import (
	"html/template"
	"math"
	"net/url"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/http"
	"github.com/gin-gonic/gin"
)

func GetTotalPage(count, pageSize int) int {
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	return totalPages
}

func PaginationForTemplate(totalPages, page int, queries url.Values) gin.H {
	queries.Del(http.PageKeyParam)

	return gin.H{
		"Page":       page,
		"TotalPages": totalPages,
		"HasPrev":    page > 1,
		"HasNext":    page < totalPages,
		"PrevPage":   page - 1,
		"NextPage":   page + 1,
		"Query":      template.URL(queries.Encode()),
	}
}
