package httperror

import (
	"net/http"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/gin-gonic/gin"
)

const (
	internalErrorRoute = "/errors/server-error"
)

func ErrorServer(server *gin.Engine, conf *config.Config) {
	templatePath := conf.APP.RootPath + conf.APP.ErrorTemplates
	server.LoadHTMLGlob(templatePath)
	server.GET(internalErrorRoute, func(ctx *gin.Context) {
		ctx.HTML(http.StatusInternalServerError, "server_error.html", gin.H{})
	})
}
