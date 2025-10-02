package router

import (
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http/handler"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/http"
	"github.com/gin-gonic/gin"
)

func meRouter(server *gin.Engine) {
	server.Use(http.AuthRequired())

	server.GET(http.PathMe, func(ctx *gin.Context) {
		handler.MeHandler(ctx)
	})
}
