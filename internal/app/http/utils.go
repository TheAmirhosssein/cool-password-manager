package http

import (
	"net/http"

	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/gin-gonic/gin"
)

func HandleError(ctx *gin.Context, customError errors.CustomError, template string, data gin.H) {
	if errors.HttpCode(customError.Code) == http.StatusInternalServerError {
		NewServerError(ctx)
	} else {
		data["error"] = true
		data["message"] = customError.Message
		ctx.HTML(errors.HttpCode(customError.Code), template, data)
		ctx.Abort()
	}
}

func HandlerFormError(ctx *gin.Context, formError error, template string, data gin.H) {
	data["error"] = true
	data["message"] = formError.Error()
	ctx.HTML(errors.HttpCode(http.StatusBadRequest), template, data)
	ctx.Abort()

}

func NewServerError(c *gin.Context) {
	c.Redirect(http.StatusFound, internalErrorRoute)
	c.Abort()
}
