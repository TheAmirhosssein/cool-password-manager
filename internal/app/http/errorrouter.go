package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	internalErrorRoute = "/errors/server-error"
)

func ErrorServer(server *gin.Engine) {
	server.GET(internalErrorRoute, func(ctx *gin.Context) {
		ctx.HTML(http.StatusInternalServerError, "server_error.html", gin.H{})
	})
}
