package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http/handler/model"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	localHttp "github.com/TheAmirhosssein/cool-password-manage/internal/app/http"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SignUpHandler(ctx *gin.Context, usecase usecase.AuthUsecase) {
	template := "sign_up.html"
	data := gin.H{"action": localHttp.PathSignUp}

	switch ctx.Request.Method {
	case http.MethodGet:
		ctx.HTML(http.StatusOK, template, data)

	case http.MethodPost:
		var form model.SignUpModel
		if err := ctx.ShouldBind(&form); err != nil {
			ctx.HTML(http.StatusOK, template, nil)
			return
		}

		acc := entity.Account{
			Username:  form.Username,
			Email:     form.Email,
			FirstName: form.FirstName,
			LastName:  form.LastName,
		}

		authenticator, err := usecase.SignUp(ctx, acc)
		if err != nil {
			localHttp.HandleError(ctx, errors.Error2Custom(err), template, data)
			return
		}

		twoFactor, err := usecase.CreateTwoFactor(ctx, acc.Username)
		if err != nil {
			localHttp.HandleError(ctx, errors.Error2Custom(err), template, data)
			return
		}

		session := sessions.Default(ctx)
		session.Set(localHttp.AuthTwoFactorIDKey, string(twoFactor.ID))

		if err := session.Save(); err != nil {
			log.ErrorLogger.Error("can not set two factor id into session", "error", err.Error(), "username", acc.Username)
			localHttp.NewServerError(ctx)
			return
		}

		base64Img := base64.StdEncoding.EncodeToString(authenticator.QrCode)

		ctx.HTML(http.StatusOK, "qrcode.html", gin.H{
			"QRCode":        base64Img,
			"twoFactorPath": localHttp.PathTwoFactor,
		})
	}
}

func LoginHandler(ctx *gin.Context, usecase usecase.AuthUsecase) {
	templateName := "login.html"
	data := gin.H{"action": localHttp.PathLogin}

	switch ctx.Request.Method {
	case http.MethodGet:
		ctx.HTML(http.StatusOK, templateName, data)

	case http.MethodPost:
		var form model.LoginModel
		if err := ctx.ShouldBind(&form); err != nil {
			formErr := errors.NewError(err.Error(), http.StatusBadRequest)
			localHttp.HandlerFormError(ctx, formErr, templateName, data)
			return
		}

		twoFactor, err := usecase.CreateTwoFactor(ctx, form.Username)
		if err != nil {
			log.ErrorLogger.Error("can not set two factor id into session", "error", err.Error(), "username", twoFactor.Username)
			localHttp.HandleError(ctx, errors.Error2Custom(err), templateName, data)
			return
		}

		session := sessions.Default(ctx)
		session.Set(localHttp.AuthTwoFactorIDKey, string(twoFactor.ID))

		if err := session.Save(); err != nil {
			localHttp.NewServerError(ctx)
			return
		}

		ctx.Redirect(http.StatusFound, localHttp.PathTwoFactor)
	}

}

func TwoFactorHandler(ctx *gin.Context, usecase usecase.AuthUsecase) {
	templateName := "two_factor.html"
	data := gin.H{"action": localHttp.PathTwoFactor}

	switch ctx.Request.Method {

	case http.MethodGet:
		ctx.HTML(http.StatusOK, templateName, data)

	case http.MethodPost:
		var form model.TwoFactorModel
		if err := ctx.ShouldBind(&form); err != nil {
			ctx.HTML(http.StatusOK, templateName, nil)
			return
		}

		session := sessions.Default(ctx)
		twoFactorID, ok := session.Get(localHttp.AuthTwoFactorIDKey).(string)

		if !ok {
			localHttp.NewServerError(ctx)
			return
		}

		account, err := usecase.ValidateTwoFactor(ctx, types.CacheID(twoFactorID), form.VerificationCode)
		if err != nil {
			localHttp.HandleError(ctx, errors.Error2Custom(err), templateName, data)
			return
		}

		session.Set(localHttp.AuthUsernameKey, account.Username)
		session.Set(localHttp.AuthUserIDKey, int64(account.Entity.ID))

		if err := session.Save(); err != nil {
			log.ErrorLogger.Error("can not set username and user id into session", "error", err.Error())
			localHttp.NewServerError(ctx)
			return
		}

		ctx.Redirect(http.StatusFound, localHttp.PathMe)
		ctx.Abort()
	}
}

func LogoutHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete(localHttp.AuthUserIDKey)
	session.Delete(localHttp.AuthUsernameKey)

	if err := session.Save(); err != nil {
		localHttp.NewServerError(ctx)
		return
	}

	ctx.Redirect(http.StatusFound, localHttp.PathLogin)
}
