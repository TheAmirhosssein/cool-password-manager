package handler

import "github.com/gin-gonic/gin"

func Authenticator(ctx *gin.Context) {
	ctx.HTML(200, "index.html", gin.H{
		"title":   "Welcome",
		"message": "Hello from Gin with HTML template!",
	})
}
