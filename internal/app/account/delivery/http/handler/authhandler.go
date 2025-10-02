package handler

import (
	"encoding/base64"
	"net/http"

	httpPath "github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http/handler/model"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/httperror"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SignUpHandler(ctx *gin.Context, usecase usecase.AuthUsecase) {
	template := "sign_up.html"
	switch ctx.Request.Method {

	case http.MethodGet:
		ctx.HTML(http.StatusOK, template, gin.H{"action": httpPath.PathSignUp})

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
			Password:  form.Password,
		}

		authenticator, err := usecase.SignUp(ctx, acc)
		if err != nil {
			httperror.HandleError(ctx, errors.Error2Custom(err), template)
			return
		}

		twoFactor, err := usecase.CreateTwoFactor(ctx, acc.Username, acc.Password)
		if err != nil {
			httperror.HandleError(ctx, errors.Error2Custom(err), template)
			return
		}

		session := sessions.Default(ctx)
		session.Set("twoFactorID", string(twoFactor.ID))

		if err := session.Save(); err != nil {
			log.ErrorLogger.Error("can not set two factor session", "error", err.Error(), "username", acc.Username)
			httperror.NewServerError(ctx)
			return
		}

		base64Img := base64.StdEncoding.EncodeToString(authenticator.QrCode)

		ctx.HTML(http.StatusOK, "qrcode.html", gin.H{
			"QRCode":        base64Img,
			"twoFactorPath": httpPath.PathTwoFactor,
		})
	}
}

func LoginHandler(ctx *gin.Context, usecase usecase.AuthUsecase) {
	templateName := "login.html"

	switch ctx.Request.Method {
	case http.MethodGet:
		ctx.HTML(http.StatusOK, templateName, gin.H{"action": httpPath.PathTwoFactor})

	case http.MethodPost:
		var form model.SignUpModel
		if err := ctx.ShouldBind(&form); err != nil {
			ctx.HTML(http.StatusOK, templateName, nil)
			return
		}

		twoFactor, err := usecase.CreateTwoFactor(ctx, form.Username, form.Password)
		if err != nil {
			httperror.HandleError(ctx, errors.Error2Custom(err), templateName)
			return
		}

		session := sessions.Default(ctx)
		session.Set("twoFactorID", string(twoFactor.ID))

		if err := session.Save(); err != nil {
			httperror.NewServerError(ctx)
			return
		}

		ctx.HTML(http.StatusOK, "qrcode.html", gin.H{
			"twoFactorID":   twoFactor.ID,
			"twoFactorPath": httpPath.PathTwoFactor,
		})
	}

}

func TwoFactorHandler(ctx *gin.Context, usecase usecase.AuthUsecase) {
	templateName := "two_factor.html"

	switch ctx.Request.Method {

	case http.MethodGet:
		ctx.HTML(http.StatusOK, templateName, gin.H{"action": httpPath.PathTwoFactor})

	case http.MethodPost:
		var form model.TwoFactorModel
		if err := ctx.ShouldBind(&form); err != nil {
			ctx.HTML(http.StatusOK, templateName, nil)
			return
		}

		session := sessions.Default(ctx)
		twoFactorID, ok := session.Get("twoFactorID").(string)
		if !ok {
			httperror.NewServerError(ctx)
			return
		}

		account, err := usecase.Login(ctx, types.CacheID(twoFactorID), form.VerificationCode)
		if err != nil {
			httperror.HandleError(ctx, errors.Error2Custom(err), templateName)
			return
		}

		session.Set("username", account.Username)

		if err := session.Save(); err != nil {
			httperror.NewServerError(ctx)
			return
		}

		ctx.Redirect(http.StatusPermanentRedirect, httpPath.PathMe)
		ctx.Abort()
	}
}
