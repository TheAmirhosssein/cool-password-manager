package repository_test

import (
	"context"
	"testing"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	params "github.com/TheAmirhosssein/cool-password-manage/internal/app/param"
	"github.com/TheAmirhosssein/cool-password-manage/internal/seed"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/base"
	"github.com/stretchr/testify/require"
)

func TestGroupRepository_Create(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewGroupRepository(pgTestSuite.db)

	testcases := []struct {
		name    string
		group   entity.Group
		wantErr bool
	}{
		{
			name: "create new group",
			group: entity.Group{
				Name:  "Cool Developers",
				Owner: seed.AccountJohnDoe,
			},
			wantErr: false,
		},
		{
			name: "duplicate group name for same owner",
			group: entity.Group{
				Name:  seed.GroupBrockhampton.Name,
				Owner: seed.GroupBrockhampton.Owner,
			},
			wantErr: true,
		},
		{
			name: "missing owner (invalid id)",
			group: entity.Group{
				Name: "Orphan Group",
				Owner: entity.Account{
					Entity: base.Entity{ID: -1},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := repo.Create(ctx, tc.group)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NoError(t, err)
			}
		})
	}
}

func TestGroupRepository_Read(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewGroupRepository(pgTestSuite.db)

	testcases := []struct {
		name    string
		param   params.ReadGroupParams
		wantErr bool
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
			wantLen: 0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			groups, err := repo.Read(ctx, tc.param)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
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

func TestGroupRepository_AddAccount(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewGroupRepository(pgTestSuite.db)

	testcases := []struct {
		name     string
		groupID  types.ID
		accounts []entity.Account
		wantErr  bool
	}{
		{
			name:    "add single account to group",
			groupID: seed.GroupBrockhampton.Entity.ID,
			accounts: []entity.Account{
				seed.AccountJohnDoe,
			},
			wantErr: false,
		},
		{
			name:    "add multiple accounts to group",
			groupID: seed.GroupBlackHippy.Entity.ID,
			accounts: []entity.Account{
				seed.AccountJohnDoe,
				seed.AccountEarl,
			},
			wantErr: false,
		},
		{
			name:    "invalid group id",
			groupID: -1,
			accounts: []entity.Account{
				seed.AccountJohnDoe,
			},
			wantErr: true,
		},
		{
			name:    "invalid account id",
			groupID: seed.GroupBrockhampton.Entity.ID,
			accounts: []entity.Account{
				{Entity: base.Entity{ID: -1}},
			},
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := repo.AddAccount(ctx, tc.groupID, tc.accounts)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
