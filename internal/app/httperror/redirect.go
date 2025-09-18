package httperror

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewServerError(c *gin.Context) {
	c.Redirect(http.StatusFound, internalErrorRoute)
	c.Abort()
}
