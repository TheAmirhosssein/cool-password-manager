package router

import (
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http/handler"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/http"
	"github.com/gin-gonic/gin"
)

func meRouter(server *gin.Engine, gRepo repository.GroupRepository, aRepo repository.AccountRepository) {
	server.Use(http.AuthRequired())
	groupeUsecase := usecase.NewGroupUsecase(gRepo, aRepo)
	server.GET(http.PathMe, func(ctx *gin.Context) {
		handler.MeHandler(ctx, groupeUsecase)
	})
}
