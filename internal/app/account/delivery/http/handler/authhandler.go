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

	ctx.HTML(http.StatusOK, template, data)
}

func SignUpInitialHandler(ctx *gin.Context, usecase usecase.AuthUsecase) {
	var body model.SignUpInitModel
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	registration := entity.Registration{
		Username:  body.Username,
		Email:     body.Email,
		FirstName: body.FirstName,
		LastName:  body.LastName,
	}

	record, cacheID, err := usecase.SignUpInit(ctx, registration, body.RegistrationRequest)
	if err != nil {
		localHttp.HandleJSONError(ctx, errors.Error2Custom(err))
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"record": base64.StdEncoding.EncodeToString(record), "registrationID": cacheID})
}

func SignUpFinalizeHandler(ctx *gin.Context, usecase usecase.AuthUsecase) {
	var body model.SignUpFinalizeModel
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	recordBytes, err := base64.StdEncoding.DecodeString(body.RegistrationRecord)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid registration record encoding"})
		return
	}

	authenticator, username, err := usecase.SignUpFinalize(ctx, recordBytes, types.CacheID(body.RegistrationID))
	if err != nil {
		localHttp.HandleJSONError(ctx, errors.Error2Custom(err))
		return
	}

	twoFactor, err := usecase.CreateTwoFactor(ctx, body.RegistrationID)
	if err != nil {
		localHttp.HandleJSONError(ctx, errors.Error2Custom(err))
		return
	}

	session := sessions.Default(ctx)
	session.Set(localHttp.AuthTwoFactorIDKey, string(twoFactor.ID))

	if err := session.Save(); err != nil {
		log.ErrorLogger.Error("can not set two factor id into session", "error", err.Error(), "username", username)
		localHttp.NewServerError(ctx)
		return
	}

	base64Img := base64.StdEncoding.EncodeToString([]byte(authenticator.QrCode))

	ctx.HTML(http.StatusOK, "qrcode.html", gin.H{
		"QRCode":        base64Img,
		"twoFactorPath": localHttp.PathTwoFactor,
	})
}

func LoginHandler(ctx *gin.Context, usecase usecase.AuthUsecase) {
	templateName := "login.html"
	data := gin.H{"signUpUrl": localHttp.PathSignUp}

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

func LoginInitHandler(ctx *gin.Context, usecase usecase.AuthUsecase) {
	var body model.LoginInitModel
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ke2, err := usecase.LoginInit(ctx, body.KE1, body.Username)
	if err != nil {
		localHttp.HandleJSONError(ctx, errors.Error2Custom(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ke2": ke2})
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
