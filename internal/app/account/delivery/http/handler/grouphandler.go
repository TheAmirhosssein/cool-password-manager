package handler

import (
	"net/http"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/param"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	localHttp "github.com/TheAmirhosssein/cool-password-manage/internal/app/http"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/convertors"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/paginator"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/gin-gonic/gin"
)

func GroupListHandler(ctx *gin.Context, usecase usecase.GroupUsecase, conf *config.Config) {
	username := ctx.GetString("username")
	templateName := "groups.html"

	page := convertors.ParseQueryParamToInt(ctx.Query(localHttp.PageKeyParam), conf.DefaultPage)
	pageSize := convertors.ParseQueryParamToInt(ctx.Query(localHttp.PageSizeKeyParam), conf.DefaultPageSize)
	limit, offset := convertors.SimplePaginationToLimitOffset(page, pageSize)

	groups, numRows, err := usecase.Read(ctx, param.ReadGroupParams{
		MemberID: types.ID(ctx.GetInt64(localHttp.AuthUserIDKey)),
		Limit:    limit,
		Offset:   offset,
	})

	if err != nil {
		localHttp.HandleError(ctx, errors.Error2Custom(err), templateName, gin.H{})
		return
	}

	ctx.HTML(http.StatusOK, templateName, gin.H{
		"username":   username,
		"logout_url": localHttp.PathLogout,
		"groups":     groups,
		"Pagination": paginator.PaginationForTemplate(paginator.GetTotalPage(numRows, pageSize), page, ctx.Request.URL.Query()),
	})
}
