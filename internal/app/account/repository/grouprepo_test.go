package repository_test

import (
	"context"
	"testing"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	params "github.com/TheAmirhosssein/cool-password-manage/internal/app/account/param"
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

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := repo.Create(ctx, &tc.group)
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

			groups, count, err := repo.Read(ctx, tc.param)
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

func TestGroupRepository_ReadOne(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewGroupRepository(pgTestSuite.db)

	testcases := []struct {
		name     string
		groupID  types.ID
		memberID types.ID
		empty    bool
	}{
		{
			name:     "valid group and valid member",
			groupID:  seed.GroupBrockhampton.Entity.ID,
			memberID: seed.AccountMattChampion.Entity.ID,
			empty:    false,
		},
		{
			name:     "valid group but member not in group",
			groupID:  seed.GroupBrockhampton.Entity.ID,
			memberID: seed.AccountEarl.Entity.ID,
			empty:    true,
		},
		{
			name:     "invalid group id",
			groupID:  -1,
			memberID: seed.AccountJohnDoe.Entity.ID,
			empty:    true,
		},
		{
			name:     "invalid member id",
			groupID:  seed.GroupBlackHippy.Entity.ID,
			memberID: -1,
			empty:    true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			group, err := repo.ReadOne(ctx, tc.groupID, tc.memberID)
			if tc.empty {
				require.Zero(t, group.Entity.ID)
				require.Empty(t, group.Name)
				require.Zero(t, group.Owner.Entity.ID)
				require.Empty(t, group.Owner.Username)
				require.Empty(t, group.Members)
			} else {
				require.NoError(t, err)
				require.NotZero(t, group.Entity.ID)
				require.NotEmpty(t, group.Name)
				require.NotZero(t, group.Owner.Entity.ID)
				require.NotEmpty(t, group.Owner.Username)
				require.NotEmpty(t, group.Members)
				for _, m := range group.Members {
					require.NotZero(t, m.Entity.ID)
					require.NotEmpty(t, m.Username)
				}
			}
		})
	}
}

func TestGroupRepository_Update(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewGroupRepository(pgTestSuite.db)

	testcases := []struct {
		name        string
		group       entity.Group
		wouldChange bool
	}{
		{
			name: "update group name successfully",
			group: entity.Group{
				Entity: base.Entity{ID: seed.GroupBrockhampton.Entity.ID},
				Name:   "Updated Brockhampton",
				Owner:  seed.GroupBrockhampton.Owner,
			},
			wouldChange: true,
		},
		{
			name: "update group name with different owner",
			group: entity.Group{
				Entity: base.Entity{ID: seed.GroupBrockhampton.Entity.ID},
				Name:   "Kendrick Brockhampton",
				Owner:  seed.AccountKendrickLamar,
			},
			wouldChange: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := repo.Update(ctx, tc.group)
			require.NoError(t, err)

			query := `SELECT name FROM groups WHERE id = $1`
			var updatedName string
			err = pgTestSuite.db.QueryRow(ctx, query, tc.group.Entity.ID).Scan(&updatedName)
			require.NoError(t, err)
			if tc.wouldChange {
				require.Equal(t, tc.group.Name, updatedName)
			} else {
				require.NotEqual(t, tc.group.Name, updatedName)
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

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := repo.AddAccounts(ctx, tc.groupID, tc.accounts)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGroupRepository_DeleteAllMembers(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewGroupRepository(pgTestSuite.db)

	testcases := []struct {
		name        string
		groupID     types.ID
		ownerID     types.ID
		wouldDelete bool
	}{
		{
			name:        "successful",
			groupID:     seed.GroupBlackHippy.Entity.ID,
			ownerID:     seed.GroupBlackHippy.Owner.Entity.ID,
			wouldDelete: true,
		},
		{
			name:        "different owner",
			groupID:     seed.GroupBrockhampton.Entity.ID,
			ownerID:     seed.GroupBlackHippy.Owner.Entity.ID,
			wouldDelete: false,
		},
	}

	for _, tc := range testcases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := repo.DeleteAllMembers(ctx, tc.groupID, tc.ownerID)
			require.NoError(t, err)

			var count int64
			query := "SELECT count(account_id) from groups_accounts WHERE groups_accounts.group_id = $1"
			pgTestSuite.db.QueryRow(ctx, query, tc.groupID).Scan(&count)

			if tc.wouldDelete {
				require.Zero(t, count)
			} else {
				require.NotZero(t, count)
			}
		})
	}
}
