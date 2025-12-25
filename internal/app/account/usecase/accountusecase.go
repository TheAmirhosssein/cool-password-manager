package usecase

import (
	"context"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
)

type AccountUsecase struct {
	accountRepo repository.AccountRepository
}

func NewAccountUsecase(accountRepo repository.AccountRepository) AccountUsecase {
	return AccountUsecase{accountRepo: accountRepo}
}

func (u *AccountUsecase) ReadByUsername(ctx context.Context, username string) (entity.Account, error) {
	exist, err := u.accountRepo.ExistByUsername(ctx, username)
	if err != nil {
		log.ErrorLogger.Error("error at checking if the user with the username exist", "error", err.Error())
		return entity.Account{}, errors.NewServerError()
	}

	if !exist {
		return entity.Account{}, account.AccountUsernameDoesNotExist
	}

	account, err := u.accountRepo.ReadByUsername(ctx, username)
	if err != nil {
		log.ErrorLogger.Error("error at reading account by username", "error", err.Error())
		return entity.Account{}, errors.NewServerError()
	}

	return account, nil
}
