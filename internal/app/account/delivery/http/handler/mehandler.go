package handler

import (
	"net/http"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	localHttp "github.com/TheAmirhosssein/cool-password-manage/internal/app/http"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/gin-gonic/gin"
)

func MeHandler(ctx *gin.Context, usecase usecase.GroupUsecase, conf *config.Config) {
	username := ctx.GetString("username")
	templateName := "me.html"

	group, err := usecase.ReadFirstGroup(ctx, types.ID(ctx.GetInt64(localHttp.AuthUserIDKey)))

	if err != nil {
		localHttp.HandleError(ctx, errors.Error2Custom(err), templateName, gin.H{})
		return
	}

	ctx.HTML(http.StatusOK, templateName, gin.H{
		"username": username, "logout_url": localHttp.PathLogout, "group": group,
	})
}
