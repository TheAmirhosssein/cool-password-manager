package usecase_test

import (
	"context"
	"testing"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	params "github.com/TheAmirhosssein/cool-password-manage/internal/app/param"
	"github.com/TheAmirhosssein/cool-password-manage/internal/seed"
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

func setupGroupUsecase() usecase.GroupUsecase {
	groupRepo := repository.NewGroupRepository(pgTestSuite.db)
	accountRepo := repository.NewAccountRepository(pgTestSuite.db)

	return usecase.NewGroupUsecase(groupRepo, accountRepo)
}
