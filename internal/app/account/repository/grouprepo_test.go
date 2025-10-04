package repository_test

import (
	"context"
	"testing"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
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
			groupID: seed.GroupRadiohead.Entity.ID,
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
