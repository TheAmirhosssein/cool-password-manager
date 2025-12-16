package http

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	AuthUserIDKey      = "user_id"
	AuthUsernameKey    = "username"
	AuthTwoFactorIDKey = "twoFactorID"
)

const (
	PageKeyParam     = "page"
	PageSizeKeyParam = "page-size"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get(AuthUsernameKey)
		userID := session.Get(AuthUserIDKey)

		if username == nil {
			c.Redirect(http.StatusFound, PathLogin)
			return
		}

		userID, ok := userID.(int64)
		if !ok {
			NewServerError(c)
			return
		}

		c.Set(AuthUserIDKey, userID)
		c.Set(AuthUsernameKey, username)

		c.Next()
	}
}

func GuestOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get(AuthUsernameKey)

		if username != nil {
			c.Redirect(http.StatusFound, PathMe)
			return
		}

		c.Next()
	}
}
