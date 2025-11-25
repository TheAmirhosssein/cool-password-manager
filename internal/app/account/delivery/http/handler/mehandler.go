package handler

import (
	"fmt"
	"net/http"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func MeHandler(ctx *gin.Context, usecase usecase.GroupUsecase) {
	userID, _ := ctx.Get("username")
	templateName := "me.html"
	session := sessions.Default(ctx)
	username := session.Get("username")
	// groups, numRows, err := usecase.Read(ctx, param.ReadGroupParams{MemberID: types.ID(userID)})
	// if err != nil {
	// 	localHttp.HandleError(ctx, errors.Error2Custom(err), templateName)
	// 	return
	// }
	fmt.Println(username)
	ctx.HTML(http.StatusOK, templateName, gin.H{"username": username, "na": userID})
}
