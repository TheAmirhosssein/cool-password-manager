package router

import (
	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func AuthHandler(server *gin.Engine, conf *config.Config) {
	templatePath := conf.APP.RootPath + conf.APP.AuthenticatorTemplates
	server.LoadHTMLGlob(templatePath)
	server.GET("/auth", handler.Authenticator)
}
