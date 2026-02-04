package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/opaque"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/totp"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/encrypt"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
)

type AuthUsecase struct {
	accountRepo      repository.AccountRepository
	twoFactorRepo    repository.TwoFactorRepository
	registrationRepo repository.RegistrationRepository

	authenticator totp.AuthenticatorAdaptor
	opaqueServer  opaque.OpaqueService

	config *config.Config
}

func NewAuthUsecase(aRepo repository.AccountRepository, tfRepo repository.TwoFactorRepository,
	rRepo repository.RegistrationRepository, authenticator totp.AuthenticatorAdaptor,
	opaqueServer opaque.OpaqueService, config *config.Config) AuthUsecase {
	return AuthUsecase{
		accountRepo:      aRepo,
		twoFactorRepo:    tfRepo,
		authenticator:    authenticator,
		opaqueServer:     opaqueServer,
		registrationRepo: rRepo,
		config:           config,
	}
}

func (u *AuthUsecase) SignUpInit(ctx context.Context, registration entity.Registration, message []byte) ([]byte, types.CacheID, error) {
	existByUsername, err := u.accountRepo.ExistByUsername(ctx, registration.Username)
	if err != nil {
		log.ErrorLogger.Error("error checking user existence by username", "error", err.Error(), "username", registration.Username)
		return nil, types.CacheID(""), errors.NewServerError()
	}

	if existByUsername {
		return nil, types.CacheID(""), account.AuthUsernameExist
	}

	existByEmail, err := u.accountRepo.ExistByEmail(ctx, registration.Email)
	if err != nil {
		log.ErrorLogger.Error("error checking user existence by email", "error", err.Error(), "username", registration.Username)
		return nil, types.CacheID(""), errors.NewServerError()
	}

	if existByEmail {
		return nil, types.CacheID(""), account.AuthEmailExist
	}

	response, credID, err := u.opaqueServer.RegisterInit(message)
	if err != nil {
		log.ErrorLogger.Error("error at registration initiation", "error", err.Error())
		return nil, types.CacheID(""), errors.NewServerError()
	}

	registration.Duration = time.Minute * time.Duration(u.config.TwoFactorDuration)
	registration.ID = types.CacheID(base64.RawURLEncoding.EncodeToString(credID))

	err = u.registrationRepo.Create(ctx, registration)
	if err != nil {
		log.ErrorLogger.Error("error at saving registration", "error", err.Error())
		return nil, types.CacheID(""), errors.NewServerError()
	}

	return response, registration.ID, nil
}

func (u *AuthUsecase) SignUpFinalize(ctx context.Context, message []byte, registrationID types.CacheID) (totp.Authenticator, string, error) {
	registration, err := u.registrationRepo.Get(ctx, registrationID)
	if err != nil {
		log.ErrorLogger.Error("error at getting registration", "error", err.Error())
		return totp.Authenticator{}, "", errors.NewServerError()
	}

	credID, err := base64.RawURLEncoding.DecodeString(string(registrationID))
	if err != nil {
		log.ErrorLogger.Error("error at converting registration id into byte", "error", err.Error())
		return totp.Authenticator{}, "", errors.NewServerError()
	}

	opaqueRecord, err := u.opaqueServer.RegisterFinalize(message, credID, registration.Username)
	if err != nil {
		log.ErrorLogger.Error("error at finalizing registration", "error", err.Error())
		return totp.Authenticator{}, "", errors.NewServerError()
	}

	acc := entity.Account{
		Username:     registration.Username,
		Email:        registration.Email,
		FirstName:    registration.FirstName,
		LastName:     registration.LastName,
		OpaqueRecord: opaqueRecord,
	}

	authenticator, err := u.authenticator.GenerateQRCode(acc.Username)
	if err != nil {
		log.ErrorLogger.Error("error at generating authenticator qr code", "error", err.Error(), "username", acc.Username)
		return totp.Authenticator{}, "", errors.NewServerError()
	}

	key, err := u.config.GetAESSecretKey()
	if err != nil {
		log.ErrorLogger.Error("error at getting aes secret key", "error", err.Error())
		return totp.Authenticator{}, "", errors.NewServerError()
	}

	secret, err := encrypt.EncryptAESSecret(key, authenticator.Secret)
	if err != nil {
		log.ErrorLogger.Error("error at encrypting authenticator secret", "error", err.Error(), "username", acc.Username)
		return totp.Authenticator{}, "", errors.NewServerError()
	}

	acc.TOTPSecret = []byte(secret)
	err = u.accountRepo.Create(ctx, acc)
	if err != nil {
		log.ErrorLogger.Error("error at creating account", "error", err.Error(), "username", acc.Username)
		return totp.Authenticator{}, "", errors.NewServerError()
	}

	return authenticator, acc.Username, nil
}

func (u *AuthUsecase) LoginInit(ctx context.Context, message []byte, username string) ([]byte, error) {
	existence, err := u.accountRepo.ExistByUsername(ctx, username)
	if err != nil {
		log.ErrorLogger.Error("error checking user existence by username", "error", err.Error(), "username", username)
		return nil, errors.NewServerError()
	}

	if !existence {
		return nil, account.AuthInvalidAccount
	}

	account, err := u.accountRepo.ReadByUsername(ctx, username)
	if err != nil {
		log.ErrorLogger.Error("error at reading user by username")
		return nil, errors.NewServerError()
	}

	message2, err := u.opaqueServer.LoginInit(message, account.OpaqueRecord, account.Username)
	if err != nil {
		log.ErrorLogger.Error("error at login initiation", "error", err.Error())
		return nil, errors.NewServerError()
	}

	return message2, nil
}

func (u *AuthUsecase) CreateTwoFactor(ctx context.Context, username string) (entity.TwoFactor, error) {
	twoFactorID, err := generateTwoFactorID()
	if err != nil {
		log.ErrorLogger.Error("error generation two factor id", "error", err.Error(), "username", username)
		return entity.TwoFactor{}, errors.NewServerError()
	}

	duration := time.Minute * time.Duration(u.config.TwoFactorDuration)
	twoFactor := entity.TwoFactor{ID: types.CacheID(twoFactorID), Username: username, Duration: duration}

	err = u.twoFactorRepo.Create(ctx, twoFactor)
	if err != nil {
		log.ErrorLogger.Error("error creating two factor", "error", err.Error(), "username", username)
		return entity.TwoFactor{}, errors.NewServerError()
	}

	log.InfoLogger.Info("two factor created for account", "username", username)
	return twoFactor, nil
}

func (u *AuthUsecase) ValidateTwoFactor(ctx context.Context, twoFactorID types.CacheID, verificationCode string) (entity.Account, error) {
	twoFactorExist, err := u.twoFactorRepo.Exist(ctx, twoFactorID)
	if err != nil {
		log.ErrorLogger.Error("error at checking if two factor exist", "error", err.Error())
		return entity.Account{}, errors.NewServerError()
	}

	if !twoFactorExist {
		return entity.Account{}, account.AuthTwoFactorDoesNotExist
	}

	twoFactor, err := u.twoFactorRepo.Get(ctx, twoFactorID)
	if err != nil {
		log.ErrorLogger.Error("error at getting two factor", "error", err.Error())
		return entity.Account{}, errors.NewServerError()
	}

	acc, err := u.accountRepo.ReadByUsername(ctx, twoFactor.Username)
	if err != nil {
		log.ErrorLogger.Error("error at reading account username", "error", err.Error(), "username", acc.Username)
		return entity.Account{}, errors.NewServerError()
	}

	key, err := u.config.GetAESSecretKey()
	if err != nil {
		log.ErrorLogger.Error("error at getting aes secret key", "error", err.Error())
		return entity.Account{}, errors.NewServerError()
	}

	secret, err := encrypt.DecryptAESSecret(key, acc.TOTPSecret)
	if err != nil {
		log.ErrorLogger.Error("error at decrypting secret", "error", err.Error(), "username", acc.Username)
		return entity.Account{}, errors.NewServerError()
	}

	codeValid := u.authenticator.VerifyCode(secret, verificationCode)
	if !codeValid {
		return entity.Account{}, account.AuthInvalidVerificationCode
	}

	return acc, nil
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
