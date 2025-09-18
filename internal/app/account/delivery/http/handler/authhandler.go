package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

func Authenticator(ctx *gin.Context) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "CoolPasswordManager",
		AccountName: "user@example.com",
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate key"})
		return
	}

	// Generate QR code PNG
	png, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create QR"})
		return
	}

	// Return base64 image for HTML <img>
	ctx.Data(http.StatusOK, "image/png", png)
}
