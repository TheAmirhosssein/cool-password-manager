package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(conf *config.Config) {
	server := gin.Default()

	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Bearer"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	config.AllowAllOrigins = true

	server.Use(cors.New(config))
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"something": "something"})
	})

	server.Run(fmt.Sprintf("%v:%v", conf.HTTP.Host, conf.HTTP.Port))
}
