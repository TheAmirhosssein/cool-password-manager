package router

import (
	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http/handler"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/totp"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func AuthRouter(server *gin.Engine, conf *config.Config, db *pgxpool.Pool, redis *redis.Client) {
	// create repository
	accountRepo := repository.NewAccountRepository(db)
	twoFactorRepo := repository.NewTwoFactorRepository(redis)
	authenticator := totp.NewAuthenticatorAdaptor(conf.Name)

	// create usecase
	authUsecase := usecase.NewAuthUsecase(accountRepo, twoFactorRepo, authenticator, conf)

	server.Any("/sign-up", func(ctx *gin.Context) {
		handler.SignUpHandler(ctx, authUsecase)
	})
}
