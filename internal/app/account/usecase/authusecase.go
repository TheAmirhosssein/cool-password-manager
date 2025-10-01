package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/totp"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/encrypt"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
)

const usecaseName = "auth"

type AuthUsecase struct {
	accountRepo   repository.AccountRepository
	twoFactorRepo repository.TwoFactorRepository
	authenticator totp.AuthenticatorAdaptor
	config        *config.Config
}

func NewAuthUsecase(aRepo repository.AccountRepository, tfRepo repository.TwoFactorRepository,
	authenticator totp.AuthenticatorAdaptor, config *config.Config) AuthUsecase {
	return AuthUsecase{accountRepo: aRepo, twoFactorRepo: tfRepo, authenticator: authenticator, config: config}
}

func (u *AuthUsecase) SignUp(ctx context.Context, acc entity.Account) (totp.Authenticator, error) {
	existByUsername, err := u.accountRepo.ExistByUsername(ctx, acc.Username)
	if err != nil {
		log.ErrorLogger.Error("error checking user existence by username", "error", err.Error(), "usecase", usecaseName)
		return totp.Authenticator{}, errors.NewServerError()
	}

	if existByUsername {
		return totp.Authenticator{}, account.AuthUsernameExist
	}

	existByEmail, err := u.accountRepo.ExistByEmail(ctx, acc.Email)
	if err != nil {
		log.ErrorLogger.Error("error checking user existence by email", "error", err.Error(), "usecase", usecaseName)
		return totp.Authenticator{}, errors.NewServerError()
	}

	if existByEmail {
		return totp.Authenticator{}, account.AuthEmailExist
	}

	authenticator, err := u.authenticator.GenerateQRCode(acc.Username)
	if err != nil {
		log.ErrorLogger.Error("error at generating authenticator qr code", "error", err.Error(), "usecase", usecaseName)
		return totp.Authenticator{}, errors.NewServerError()
	}

	secret, err := encrypt.EncryptAESSecret([]byte(u.config.GetAESSecretKey()), authenticator.Secret)
	if err != nil {
		log.ErrorLogger.Error("error at encrypting authenticator secret", "error", err.Error(), "usecase", usecaseName)
		return totp.Authenticator{}, errors.NewServerError()
	}

	password, err := encrypt.HashPassword(acc.Password)
	if err != nil {
		log.ErrorLogger.Error("error at hashing password", "error", err.Error(), "usecase", usecaseName)
		return totp.Authenticator{}, errors.NewServerError()
	}

	acc.Secret = secret
	acc.Password = password
	err = u.accountRepo.Create(ctx, acc)
	if err != nil {
		log.ErrorLogger.Error("error at creating account", "error", err.Error(), "usecase", usecaseName)
		return totp.Authenticator{}, errors.NewServerError()
	}

	return authenticator, nil
}

func (u *AuthUsecase) CreateTwoFactor(ctx context.Context, username, password string) (entity.TwoFactor, error) {
	existence, err := u.accountRepo.ExistByUsername(ctx, username)
	if err != nil {
		log.ErrorLogger.Error("error checking user existence by username", "error", err.Error(), "usecase", usecaseName)
		return entity.TwoFactor{}, errors.NewServerError()
	}

	if !existence {
		return entity.TwoFactor{}, account.AuthInvalidAccount
	}

	acc, err := u.accountRepo.ReadByUsername(ctx, username)
	if err != nil {
		log.ErrorLogger.Error("error getting account by username", "error", err.Error(), "usecase", usecaseName)
		return entity.TwoFactor{}, errors.NewServerError()
	}

	correctPassword := encrypt.CheckPasswordHash(password, acc.Password)
	if !correctPassword {
		return entity.TwoFactor{}, account.AuthInvalidAccount
	}

	twoFactorID, err := generateTwoFactorID()
	if err != nil {
		log.ErrorLogger.Error("error generation two factor id", "error", err.Error(), "usecase", usecaseName)
		return entity.TwoFactor{}, errors.NewServerError()
	}

	duration := time.Minute * 2
	twoFactor := entity.TwoFactor{ID: types.CacheID(twoFactorID), Username: username, Duration: duration}

	err = u.twoFactorRepo.Create(ctx, twoFactor)
	if err != nil {
		log.ErrorLogger.Error("error creating two factor", "error", err.Error(), "usecase", usecaseName)
		return entity.TwoFactor{}, errors.NewServerError()
	}

	log.InfoLogger.Info("two factor created for account", "username", username)
	return twoFactor, nil
}

func generateTwoFactorID() (string, error) {
	characterLength := 16
	bytes := make([]byte, characterLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil

}
