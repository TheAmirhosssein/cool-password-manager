package router

import (
	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http/handler"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/http"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/opaque"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/totp"
	"github.com/gin-gonic/gin"
)

func authRouter(
	server *gin.Engine, aRepo repository.AccountRepository, tfRepo repository.TwoFactorRepository, rRepo repository.RegistrationRepository,
	totp totp.AuthenticatorAdaptor, opaqueAdaptor opaque.OpaqueService, conf *config.Config,
) {
	authUsecase := usecase.NewAuthUsecase(aRepo, tfRepo, rRepo, totp, opaqueAdaptor, conf)

	server.GET(http.PathSignUp, http.GuestOnly(), func(ctx *gin.Context) {
		handler.SignUpHandler(ctx, authUsecase)
	})

	server.POST(http.PathSignUpInit, http.GuestOnly(), func(ctx *gin.Context) {
		handler.SignUpInitialHandler(ctx, authUsecase)
	})

	server.POST(http.PathSignUpFinal, http.GuestOnly(), func(ctx *gin.Context) {
		handler.SignUpFinalizeHandler(ctx, authUsecase)
	})

	server.GET(http.PathLogin, http.GuestOnly(), func(ctx *gin.Context) {
		handler.LoginHandler(ctx, authUsecase)
	})

	server.POST(http.PathLoginInit, http.GuestOnly(), func(ctx *gin.Context) {
		handler.LoginInitHandler(ctx, authUsecase)
	})

	server.GET(http.PathTwoFactor, http.GuestOnly(), func(ctx *gin.Context) {
		handler.TwoFactorHandler(ctx, authUsecase)
	})

	server.POST(http.PathTwoFactor, http.GuestOnly(), func(ctx *gin.Context) {
		handler.TwoFactorHandler(ctx, authUsecase)
	})

	server.GET(http.PathLogout, handler.LogoutHandler)
}
