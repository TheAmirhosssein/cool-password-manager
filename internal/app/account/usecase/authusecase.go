package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/hash"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
)

type AuthUsecase struct {
	accountRepo   repository.AccountRepository
	twoFactorRepo repository.TwoFactorRepository
}

func NewAuthUsecase(aRepo repository.AccountRepository, tfRepo repository.TwoFactorRepository) AuthUsecase {
	return AuthUsecase{accountRepo: aRepo, twoFactorRepo: tfRepo}
}

func (u *AuthUsecase) CreateTwoFactor(ctx context.Context, username, password string) (entity.TwoFactor, error) {
	existence, err := u.accountRepo.ExistByUsername(ctx, username)
	if err != nil {
		log.ErrorLogger.Error("error checking user existence by username", "error", err.Error(), "username", username)
		return entity.TwoFactor{}, errors.NewServerError()
	}

	if !existence {
		return entity.TwoFactor{}, account.AuthInvalidUser
	}

	acc, err := u.accountRepo.ReadByUsername(ctx, username)
	if err != nil {
		log.ErrorLogger.Error("error getting account by username", "error", err.Error(), "username", username)
		return entity.TwoFactor{}, errors.NewServerError()
	}

	correctPassword := hash.CheckPasswordHash(password, acc.Password)
	if !correctPassword {
		return entity.TwoFactor{}, account.AuthInvalidUser
	}

	twoFactorID, err := generateTwoFactorID()
	if err != nil {
		log.ErrorLogger.Error("error generation two factor id", "error", err.Error(), "username", username)
		return entity.TwoFactor{}, errors.NewServerError()
	}

	duration := time.Minute * 2
	twoFactor := entity.TwoFactor{ID: types.CacheID(twoFactorID), Username: username, Duration: duration}

	err = u.twoFactorRepo.Create(ctx, twoFactor)
	if err != nil {
		log.ErrorLogger.Error("error creating two factor", "error", err.Error(), "username", username)
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
