package totp

import (
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

type Authenticator struct {
	Barcode []byte
	Secret  string
	URL     string
}

type AuthenticatorAdaptor struct {
	Issuer string
}

func NewAuthenticatorAdaptor(issuer string) AuthenticatorAdaptor {
	return AuthenticatorAdaptor{Issuer: issuer}
}

func (a *AuthenticatorAdaptor) GenerateQRCode(accountName string) (Authenticator, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      a.Issuer,
		AccountName: accountName,
	})
	if err != nil {
		log.ErrorLogger.Error("error at generation authenticator qr code", "error", err.Error(), "account_name", accountName)
		return Authenticator{}, err
	}

	png, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		log.ErrorLogger.Error("error at encoding qr code to png", "error", err.Error(), "account_name", accountName)
		return Authenticator{}, err
	}

	return Authenticator{Barcode: png, Secret: key.Secret(), URL: key.URL()}, nil
}

func (a *AuthenticatorAdaptor) VerifyCode(secret, code string) bool {
	return totp.Validate(code, secret)
}
