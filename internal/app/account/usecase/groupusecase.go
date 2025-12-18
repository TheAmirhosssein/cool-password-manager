package usecase

import (
	"context"
	"slices"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	params "github.com/TheAmirhosssein/cool-password-manage/internal/app/account/param"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
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

func (u GroupUsecase) Read(ctx context.Context, params params.ReadGroupParams) ([]entity.Group, int64, error) {
	return u.groupRepo.Read(ctx, params)
}

func (u GroupUsecase) ReadFirstGroup(ctx context.Context, memberID types.ID) (entity.Group, error) {
	params := params.ReadGroupParams{MemberID: memberID, Limit: 1, Offset: 0}
	groups, _, err := u.groupRepo.Read(ctx, params)
	if err != nil {
		log.ErrorLogger.Error("error at reading last group", "error", err.Error())
		return entity.Group{}, errors.NewServerError()
	}

	if len(groups) == 0 {
		return entity.Group{}, nil
	}

	return groups[0], nil
}

func (u GroupUsecase) Update(ctx context.Context, editorAccount entity.Account, group entity.Group) error {
	toBeUpdatedGroup, err := u.groupRepo.ReadOne(ctx, group.ID, editorAccount.Entity.ID)
	if err != nil {
		log.ErrorLogger.Error("error at getting group by id", "error", err.Error())
		return errors.NewServerError()
	}

	if editorAccount.Entity.ID != toBeUpdatedGroup.Owner.Entity.ID {
		return account.GroupOnlyTheOwnerCanEdit
	}

	err = u.groupRepo.Update(ctx, group)
	if err != nil {
		log.ErrorLogger.Error("error at updating group", "error", err.Error())
		return errors.NewServerError()
	}

	if !slices.Contains(group.Members, group.Owner) {
		group.Members = append(group.Members, entity.Account{Entity: group.Owner.Entity})
	}

	err = u.groupRepo.DeleteAllMembers(ctx, group.ID, group.Owner.Entity.ID)
	if err != nil {
		log.ErrorLogger.Error("error at deleting all the members of group", "error", err.Error())
		return errors.NewServerError()
	}

	err = u.groupRepo.AddAccounts(ctx, group.ID, group.Members)
	if err != nil {
		log.ErrorLogger.Error("error at adding members to group", "error", err.Error())
		return errors.NewServerError()
	}

	return nil
}
