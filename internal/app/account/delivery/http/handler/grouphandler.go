package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http/handler/model"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/param"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	localHttp "github.com/TheAmirhosssein/cool-password-manage/internal/app/http"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/base"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/convertors"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/paginator"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/gin-gonic/gin"
)

func GroupListHandler(ctx *gin.Context, usecase usecase.GroupUsecase, conf *config.Config) {
	username := ctx.GetString(localHttp.AuthUsernameKey)
	templateName := "groups.html"

	page := convertors.ParseQueryParamToInt(ctx.Query(localHttp.PageKeyParam), conf.DefaultPage)
	pageSize := convertors.ParseQueryParamToInt(ctx.Query(localHttp.PageSizeKeyParam), conf.DefaultPageSize)
	limit, offset := convertors.SimplePaginationToLimitOffset(page, pageSize)
	searchQuery := ctx.Query(localHttp.SearchKeyParam)

	groups, numRows, err := usecase.Read(ctx, param.ReadGroupParams{
		MemberID:    types.ID(ctx.GetInt64(localHttp.AuthUserIDKey)),
		SearchQuery: types.NewNullString(searchQuery),
		Limit:       limit,
		Offset:      offset,
	})

	if err != nil {
		localHttp.HandleError(ctx, errors.Error2Custom(err), templateName, gin.H{})
		return
	}

	ctx.HTML(http.StatusOK, templateName, gin.H{
		"Username":    username,
		"LogoutUrl":   localHttp.PathLogout,
		"EditPath":    localHttp.PathGroupEdit,
		"DeletePath":  localHttp.PathGroupDelete,
		"Groups":      groups,
		"Pagination":  paginator.PaginationForTemplate(paginator.GetTotalPage(numRows, pageSize), page, ctx.Request.URL.Query()),
		"SearchQuery": searchQuery,
		"CreateUrl":   localHttp.PathGroupCreate,
	})
}

func GroupCreateHandler(ctx *gin.Context, usecase usecase.GroupUsecase) {
	templateName := "group_create.html"
	userID := types.ID(ctx.GetInt64(localHttp.AuthUserIDKey))
	data := gin.H{
		"SearchUrl": localHttp.PathGroupSearchMember,
		"Action":    localHttp.PathGroupCreate,
		"LogoutUrl": localHttp.PathLogout,
		"Username":  ctx.GetString(localHttp.AuthUsernameKey),
	}

	switch ctx.Request.Method {
	case http.MethodGet:
		ctx.HTML(http.StatusOK, templateName, data)

	case http.MethodPost:
		var form model.GroupCreate
		if err := ctx.ShouldBind(&form); err != nil {
			formErr := errors.NewError(err.Error(), http.StatusBadRequest)
			localHttp.HandlerFormError(ctx, formErr, templateName, data)
			return
		}

		group := entity.Group{
			Name:        form.Name,
			Description: types.NewNullString(form.Description),
			Owner:       entity.Account{Entity: base.Entity{ID: userID}},
		}

		for _, memberID := range form.MembersID {
			group.Members = append(group.Members, entity.Account{Entity: base.Entity{ID: memberID}})
		}

		err := usecase.Create(ctx, &group)
		if err != nil {
			localHttp.HandleError(ctx, errors.Error2Custom(err), templateName, data)
			return
		}

		ctx.Redirect(http.StatusSeeOther, localHttp.PathGroupList)
		return
	}
}

func GroupEditHandler(ctx *gin.Context, usecase usecase.GroupUsecase) {
	templateName := "group_edit.html"
	userID := types.ID(ctx.GetInt64(localHttp.AuthUserIDKey))
	groupID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		localHttp.HandleNotFoundError(ctx)
		return
	}

	data := gin.H{
		"SearchUrl": localHttp.PathGroupSearchMember,
		"Action":    fmt.Sprint(localHttp.PathGroupEdit, groupID),
		"LogoutUrl": localHttp.PathLogout,
		"Username":  ctx.GetString(localHttp.AuthUsernameKey),
	}

	switch ctx.Request.Method {
	case http.MethodGet:
		group, err := usecase.ReadOne(ctx, types.ID(groupID), userID)
		if err != nil {
			localHttp.HandleError(ctx, errors.Error2Custom(err), templateName, data)
			return
		}
		data["Group"] = group
		ctx.HTML(http.StatusOK, templateName, data)

	case http.MethodPost:
		var form model.GroupUpdate
		if err := ctx.ShouldBind(&form); err != nil {
			formErr := errors.NewError(err.Error(), http.StatusBadRequest)
			localHttp.HandlerFormError(ctx, formErr, templateName, data)
			return
		}

		group := entity.Group{
			Entity:      base.Entity{ID: types.ID(groupID)},
			Name:        form.Name,
			Description: types.NewNullString(form.Description),
			Owner:       entity.Account{Entity: base.Entity{ID: userID}},
		}

		for _, memberID := range form.MembersID {
			group.Members = append(group.Members, entity.Account{Entity: base.Entity{ID: memberID}})
		}

		err = usecase.Update(ctx, group.Owner, group)
		if err != nil {
			localHttp.HandleError(ctx, errors.Error2Custom(err), templateName, data)
			return
		}

		ctx.Redirect(http.StatusSeeOther, localHttp.PathGroupList)
	}
}

func GroupDeleteHandler(ctx *gin.Context, usecase usecase.GroupUsecase) {
	userID := types.ID(ctx.GetInt64(localHttp.AuthUserIDKey))
	groupID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		localHttp.HandleNotFoundError(ctx)
		return
	}

	err = usecase.Delete(ctx, types.ID(groupID), userID)
	if err != nil {
		localHttp.HandleError(ctx, errors.Error2Custom(err), "general_error.html", gin.H{})
		return
	}

	ctx.Redirect(http.StatusSeeOther, localHttp.PathGroupList)
}

func GroupSearchMember(ctx *gin.Context, usecase usecase.GroupUsecase) {
	username := ctx.Query("username")
	account, err := usecase.SearchMember(ctx, username)
	if err != nil {
		localHttp.HandleJSONError(ctx, errors.Error2Custom(err))
		return
	}

	option := fmt.Sprintf(`<option value="%d" selected>%s</option>`, account.Entity.ID, account.Username)
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(http.StatusOK, option)
}
