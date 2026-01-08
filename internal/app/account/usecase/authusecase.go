package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

	registration.CredID = credID
	registration.Duration = time.Minute * time.Duration(u.config.TwoFactorDuration)
	registration.ID, err = generateRegistrationID(registration.Username)
	if err != nil {
		log.ErrorLogger.Error("error at generation registration id", "error", err.Error())
		return nil, types.CacheID(""), errors.NewServerError()
	}

	err = u.registrationRepo.Create(ctx, registration)
	if err != nil {
		log.ErrorLogger.Error("error at saving registration", "error", err.Error())
		return nil, types.CacheID(""), errors.NewServerError()
	}

	return response, registration.ID, nil
}

func (u *AuthUsecase) SignUp(ctx context.Context, acc entity.Account) (totp.Authenticator, error) {
	authenticator, err := u.authenticator.GenerateQRCode(acc.Username)
	if err != nil {
		log.ErrorLogger.Error("error at generating authenticator qr code", "error", err.Error(), "username", acc.Username)
		return totp.Authenticator{}, errors.NewServerError()
	}

	key, err := u.config.GetAESSecretKey()
	if err != nil {
		log.ErrorLogger.Error("error at getting aes secret key", "error", err.Error())
		return totp.Authenticator{}, errors.NewServerError()
	}

	secret, err := encrypt.EncryptAESSecret(key, authenticator.Secret)
	if err != nil {
		log.ErrorLogger.Error("error at encrypting authenticator secret", "error", err.Error(), "username", acc.Username)
		return totp.Authenticator{}, errors.NewServerError()
	}

	acc.TOTPSecret = []byte(secret)
	err = u.accountRepo.Create(ctx, acc)
	if err != nil {
		log.ErrorLogger.Error("error at creating account", "error", err.Error(), "username", acc.Username)
		return totp.Authenticator{}, errors.NewServerError()
	}

	return authenticator, nil
}

func (u *AuthUsecase) CreateTwoFactor(ctx context.Context, username string) (entity.TwoFactor, error) {
	existence, err := u.accountRepo.ExistByUsername(ctx, username)
	if err != nil {
		log.ErrorLogger.Error("error checking user existence by username", "error", err.Error(), "username", username)
		return entity.TwoFactor{}, errors.NewServerError()
	}

	if !existence {
		return entity.TwoFactor{}, account.AuthInvalidAccount
	}

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

func generateRegistrationID(username string) (types.CacheID, error) {
	characterLength := 16
	bytes := make([]byte, characterLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return types.CacheID(fmt.Sprintf(hex.EncodeToString(bytes), username)), nil
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
