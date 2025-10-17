package usecase

import (
	"context"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
)

type GroupUsecase struct {
	groupRepo   repository.GroupRepository
	accountRepo repository.AccountRepository
}

func NewGroupUsecase(groupRepo repository.GroupRepository, accountRepo repository.AccountRepository) GroupUsecase {
	return GroupUsecase{groupRepo: groupRepo, accountRepo: accountRepo}
}

func (u GroupUsecase) Create(ctx context.Context, group *entity.Group) error {
	err := u.groupRepo.Create(ctx, group)
	if err != nil {
		log.ErrorLogger.Error("error at creating group", "error", err.Error())
		return errors.NewServerError()
	}

	group.Members = append(group.Members, entity.Account{Entity: group.Owner.Entity})

	err = u.groupRepo.AddAccounts(ctx, group.ID, group.Members)
	if err != nil {
		log.ErrorLogger.Error("error at adding members into group", "error", err.Error())
		return errors.NewServerError()
	}

	return nil
}
