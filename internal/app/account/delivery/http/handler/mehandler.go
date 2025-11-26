package handler

import (
	"net/http"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	localHttp "github.com/TheAmirhosssein/cool-password-manage/internal/app/http"
	"github.com/gin-gonic/gin"
)

func MeHandler(ctx *gin.Context, usecase usecase.GroupUsecase) {
	username := ctx.GetString("username")
	templateName := "me.html"

	ctx.HTML(http.StatusOK, templateName, gin.H{"username": username, "logout_url": localHttp.PathLogout})
}
