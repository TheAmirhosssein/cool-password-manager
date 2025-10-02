package http

import (
	"net/http"

	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/gin-gonic/gin"
)

func HandleError(ctx *gin.Context, customError errors.CustomError, template string) {
	if errors.HttpCode(customError.Code) == http.StatusInternalServerError {
		NewServerError(ctx)
	} else {
		ctx.HTML(errors.HttpCode(customError.Code), template, gin.H{"error": true, "message": customError.Message})
		ctx.Abort()
	}
}

func NewServerError(c *gin.Context) {
	c.Redirect(http.StatusFound, internalErrorRoute)
	c.Abort()
}
