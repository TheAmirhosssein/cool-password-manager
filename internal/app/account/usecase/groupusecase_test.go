package usecase_test

import (
	"context"
	"testing"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	params "github.com/TheAmirhosssein/cool-password-manage/internal/app/account/param"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	"github.com/TheAmirhosssein/cool-password-manage/internal/seed"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/base"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestGroupUsecase_Create(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	u := setupGroupUsecase()

	johnDoe := seed.AccountJohnDoe

	testcases := []struct {
		name        string
		group       entity.Group
		expectedErr error
	}{
		{
			name: "success",
			group: entity.Group{
				Name: "John's Secure Group",
				Owner: entity.Account{
					Entity: base.Entity{ID: johnDoe.Entity.ID},
				},
			},
			expectedErr: nil,
		},
		{
			name: "invalid owner id (does not exist)",
			group: entity.Group{
				Name: "Ghost Group",
				Owner: entity.Account{
					Entity: base.Entity{ID: 999999},
				},
			},
			expectedErr: errors.NewServerError(),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := u.Create(ctx, &tc.group)

			if tc.expectedErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				groupRepo := repository.NewGroupRepository(pgTestSuite.db)
				createdGroup, err := groupRepo.ReadOne(ctx, tc.group.ID, tc.group.Owner.Entity.ID)
				require.NoError(t, err)
				require.Equal(t, tc.group.Name, createdGroup.Name)
				require.NotEmpty(t, createdGroup.Members)
				require.Equal(t, johnDoe.Entity.ID, createdGroup.Members[0].Entity.ID)
			}
		})
	}
}

func TestGroupUsecase_Read(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	usecase := setupGroupUsecase()

	testcases := []struct {
		name    string
		param   params.ReadGroupParams
		wantErr bool
		count   int64
		wantLen int
	}{
		{
			name: "owner has groups with members",
			param: params.ReadGroupParams{
				MemberID: seed.AccountMattChampion.Entity.ID,
				Limit:    10,
				Offset:   0,
			},
			wantErr: false,
			count:   1,
			wantLen: 1,
		},
		{
			name: "owner has no groups",
			param: params.ReadGroupParams{
				MemberID: -1,
				Limit:    10,
				Offset:   0,
			},
			wantErr: false,
			count:   0,
			wantLen: 0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			groups, count, err := usecase.Read(ctx, tc.param)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, count, tc.count)
				for _, g := range groups {
					require.NotZero(t, g.Entity.ID)
					require.NotEmpty(t, g.Name)
					require.NotZero(t, g.Owner.Entity.ID)
					require.NotEmpty(t, g.Owner.Username)
					require.NotEmpty(t, g.Members)
				}
			}
		})
	}
}

func TestGroupUsecase_ReadFirstGroup(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	usecase := setupGroupUsecase()

	testcases := []struct {
		name      string
		accountID types.ID
		lastGroup entity.Group
		wantErr   bool
	}{
		{
			name:      "successful",
			accountID: seed.AccountEarl.Entity.ID,
			lastGroup: seed.GroupOddFuture,
			wantErr:   false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			group, err := usecase.ReadFirstGroup(ctx, tc.accountID)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, group.Name, tc.lastGroup.Name)
			}
		})
	}
}

func TestGroupUsecase_ReadOne(t *testing.T) {
	t.Parallel()
	expectedGroup := seed.GroupOddFuture
	ctx := context.Background()
	usecase := setupGroupUsecase()

	testcases := []struct {
		name          string
		groupID       types.ID
		accountID     types.ID
		expectedGroup entity.Group
		wantErr       bool
	}{
		{
			name:          "successful",
			groupID:       expectedGroup.ID,
			accountID:     seed.AccountEarl.Entity.ID,
			expectedGroup: expectedGroup,
			wantErr:       false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			group, err := usecase.ReadOne(ctx, tc.groupID, tc.accountID)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, group.Name, tc.expectedGroup.Name)
			}
		})
	}
}

func TestGroupUsecase_Update(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	usecase := setupGroupUsecase()

	g := seed.GroupBlackHippy

	testcases := []struct {
		name  string
		group entity.Group
		err   error
	}{
		{
			name: "success",
			group: entity.Group{
				Entity:      base.Entity{ID: g.ID},
				Name:        "new group name",
				Description: types.NullString{String: "something new", Valid: true},
				Owner:       g.Owner,
				Members:     []entity.Account{seed.AccountEarl, seed.AccountFrankOcean, seed.AccountKendrickLamar},
			},
		},
		{
			name: "different owner",
			group: entity.Group{
				Entity:      base.Entity{ID: g.ID},
				Name:        "new group",
				Description: types.NullString{String: "something new", Valid: true},
				Owner:       seed.GroupBrockhampton.Owner,
				Members:     []entity.Account{seed.AccountEarl, seed.AccountFrankOcean},
			},
			err: account.GroupOnlyTheOwnerCanEdit,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := usecase.Update(ctx, tc.group.Owner, tc.group)
			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err, tc.err)
				groupRepo := repository.NewGroupRepository(pgTestSuite.db)
				group, err := groupRepo.ReadOne(ctx, tc.group.ID, tc.group.Owner.Entity.ID)
				require.NoError(t, err)
				require.Equal(t, group.ID, tc.group.ID)
				require.Equal(t, group.Name, tc.group.Name)
				require.Equal(t, group.Description, tc.group.Description)
				require.Equal(t, group.Owner.Entity.ID, tc.group.Owner.Entity.ID)
				for i, member := range group.Members {
					require.Equal(t, member.Entity.ID, tc.group.Members[i].Entity.ID)
				}
			}
		})
	}
}

func setupGroupUsecase() usecase.GroupUsecase {
	groupRepo := repository.NewGroupRepository(pgTestSuite.db)
	accountRepo := repository.NewAccountRepository(pgTestSuite.db)

	return usecase.NewGroupUsecase(groupRepo, accountRepo)
}
