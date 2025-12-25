package usecase_test

import (
	"context"
	"testing"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	"github.com/TheAmirhosssein/cool-password-manage/internal/seed"
	"github.com/stretchr/testify/require"
)

func TestAccountRepository_ReadByUsername(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	usecase := setupAccountUsecase()

	testcases := []struct {
		name     string
		username string
		expect   entity.Account
		wantErr  bool
	}{
		{
			name:     "existing user",
			username: seed.AccountJohnDoe.Username,
			expect:   seed.AccountJohnDoe,
			wantErr:  false,
		},
		{
			name:     "non-existing user",
			username: "not_found",
			expect:   entity.Account{},
			wantErr:  true,
		},
	}

	for _, tc := range testcases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			account, err := usecase.ReadByUsername(ctx, tc.username)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expect.Username, account.Username)
				require.Equal(t, tc.expect.Email, account.Email)
				require.Equal(t, tc.expect.Secret, account.Secret)
			}
		})
	}
}

func setupAccountUsecase() usecase.AccountUsecase {
	aRepo := repository.NewAccountRepository(pgTestSuite.db)

	return usecase.NewAccountUsecase(aRepo)
}
