package router

import (
	"fmt"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http/handler"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/http"
	"github.com/gin-gonic/gin"
)

func groupRouter(server *gin.Engine, gRepo repository.GroupRepository, aRepo repository.AccountRepository, conf *config.Config) {
	server.Use(http.AuthRequired())
	groupeUsecase := usecase.NewGroupUsecase(gRepo, aRepo)
	server.GET(http.PathGroupList, func(ctx *gin.Context) {
		handler.GroupListHandler(ctx, groupeUsecase, conf)
	})
	server.GET(fmt.Sprint(http.PathGroupEdit, ":id/"), func(ctx *gin.Context) {
		handler.GroupEditHandler(ctx, groupeUsecase, conf)
	})
	server.POST(fmt.Sprint(http.PathGroupEdit, ":id/"), func(ctx *gin.Context) {
		handler.GroupEditHandler(ctx, groupeUsecase, conf)
	})
	server.GET(fmt.Sprint(http.PathGroupSearchMember, ":username/"), func(ctx *gin.Context) {
		handler.GroupSearchMember(ctx, groupeUsecase)
	})
}
