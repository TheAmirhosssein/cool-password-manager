package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func MeHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	user := session.Get("username")
	ctx.HTML(http.StatusOK, "me.html", gin.H{"username": user})
}
