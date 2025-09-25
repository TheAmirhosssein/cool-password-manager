package totp

import (
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

type AuthenticatorAdaptor struct {
	Issuer string
}

func (a *AuthenticatorAdaptor) GenerateQRCode(accountName string) ([]byte, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      a.Issuer,
		AccountName: accountName,
	})
	if err != nil {
		log.ErrorLogger.Error("error at generation authenticator qr code", "error", err.Error(), "account_name", accountName)
		return nil, "", err
	}

	png, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		log.ErrorLogger.Error("error at encoding qr code to png", "error", err.Error(), "account_name", accountName)
		return nil, "", err
	}

	return png, key.Secret(), nil
}
