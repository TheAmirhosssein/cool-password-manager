package router

import (
	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/totp"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func AccountRouter(server *gin.Engine, conf *config.Config, db *pgxpool.Pool, redis *redis.Client) {
	store := cookie.NewStore([]byte(conf.SecretKey))
	server.Use(sessions.Sessions("mysession", store))

	// Create auth
	accountRepo := repository.NewAccountRepository(db)
	twoFactorRepo := repository.NewTwoFactorRepository(redis)
	groupRepo := repository.NewGroupRepository(db)
	authenticator := totp.NewAuthenticatorAdaptor(conf.Name)

	// Register routers
	authRouter(server, accountRepo, twoFactorRepo, authenticator, conf)
	meRouter(server, groupRepo, accountRepo, conf)
	groupRouter(server, groupRepo, accountRepo, conf)
}
