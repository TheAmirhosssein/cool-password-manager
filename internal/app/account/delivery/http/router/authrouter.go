package router

import (
	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http/handler"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/totp"
	"github.com/gin-gonic/gin"
)

func authRouter(
	server *gin.Engine, aRepo repository.AccountRepository, tfRepo repository.TwoFactorRepository, totp totp.AuthenticatorAdaptor,
	conf *config.Config,
) {
	authUsecase := usecase.NewAuthUsecase(aRepo, tfRepo, totp, conf)

	server.Any(http.PathSignUp, func(ctx *gin.Context) {
		handler.SignUpHandler(ctx, authUsecase)
	})

	server.Any(http.PathLogin, func(ctx *gin.Context) {
		handler.LoginHandler(ctx, authUsecase)
	})

	server.Any(http.PathTwoFactor, func(ctx *gin.Context) {
		handler.TwoFactorHandler(ctx, authUsecase)
	})
}
