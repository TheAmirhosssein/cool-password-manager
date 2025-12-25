package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/delivery/http/router"
	localHttp "github.com/TheAmirhosssein/cool-password-manage/internal/app/http"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/database"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context, conf *config.Config) error {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Bearer"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}))

	db := database.GetDb(ctx)
	redisClient := redis.GetClient()

	server.LoadHTMLGlob(conf.APP.RootPath + conf.APP.TemplatePath)
	server.Static(conf.APP.StaticPath, conf.APP.RootPath+conf.APP.StaticPath)

	router.AccountRouter(server, conf, db, redisClient)
	localHttp.ErrorServer(server)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", conf.HTTP.Host, conf.HTTP.Port),
		Handler: server,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
