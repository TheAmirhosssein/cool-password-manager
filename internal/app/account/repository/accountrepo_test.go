package repository_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/database"
	"github.com/TheAmirhosssein/cool-password-manage/internal/seed"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/testdocker"
	"github.com/alicebob/miniredis/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type postgresTest struct {
	db   *pgxpool.Pool
	name string
}

var pgTestSuite postgresTest
var redisClient *redis.Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgTask := make(chan postgresTest)
	go func() {
		pgName, pgTest := database.SetupTestDB(ctx)
		pgTask <- postgresTest{name: pgName, db: pgTest}
	}()

	pgTestSuite = <-pgTask

	seed.CreateSeed(ctx, pgTestSuite.db)

	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("An error occurred while starting miniredis: %v", err)
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	seed.CreateTowFactorSeed(ctx, redisClient)

	exitCode := m.Run()
	testdocker.StopAndRemoveContainer(ctx, pgTestSuite.name, pgTestSuite.name)

	os.Exit(exitCode)
}

func TestAccountRepository_Create(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewAccountRepository(pgTestSuite.db)

	testcases := []struct {
		name    string
		account entity.Account
		wantErr bool
	}{
		{
			name: "create new account",
			account: entity.Account{
				Username:  "new_user",
				Password:  "secure_password",
				Email:     "new_user@example.com",
				FirstName: "New",
				LastName:  "User",
				Secret:    "new_secret",
			},
			wantErr: false,
		},
		{
			name: "duplicate username",
			account: entity.Account{
				Username:  seed.AccountJohnDoe.Username,
				Password:  "somepass",
				Email:     "duplicate@example.com",
				FirstName: "Dup",
				LastName:  "User",
				Secret:    "new_secret",
			},
			wantErr: true,
		},
		{
			name: "duplicate email",
			account: entity.Account{
				Username:  "another_user",
				Password:  "somepass",
				Email:     seed.AccountJohnDoe.Email,
				FirstName: "Another",
				LastName:  "User",
				Secret:    "new_secret",
			},
			wantErr: true,
		},
	}

	for _, tc := range testcases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := repo.Create(ctx, tc.account)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				account, err := repo.ReadByUsername(ctx, tc.account.Username)
				require.NoError(t, err)
				require.Equal(t, tc.account.Username, account.Username)
				require.Equal(t, tc.account.Email, account.Email)
			}
		})
	}
}

func TestAccountRepository_ReadByUsername(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewAccountRepository(pgTestSuite.db)

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

			account, err := repo.ReadByUsername(ctx, tc.username)
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

func TestAccountRepository_Update(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewAccountRepository(pgTestSuite.db)

	account := seed.AccountJohnDoe
	account.Secret = "something else"

	testcases := []struct {
		name    string
		account entity.Account
	}{
		{
			name:    "valid update",
			account: account,
		},
	}

	for _, tc := range testcases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := repo.Update(ctx, tc.account)
			require.NoError(t, err)
			account, err = repo.ReadByUsername(ctx, tc.account.Username)
			require.NoError(t, err)
			require.Equal(t, account.Secret, tc.account.Secret)
		})
	}
}

func TestAccountRepository_ExistByUsername(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewAccountRepository(pgTestSuite.db)

	testcases := []struct {
		name     string
		username string
		exist    bool
	}{
		{
			name:     "exist",
			username: seed.AccountJohnDoe.Username,
			exist:    true,
		},
		{
			name:     "wrong username",
			username: "wrong username",
			exist:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			exist, err := repo.ExistByUsername(ctx, tc.username)
			require.NoError(t, err)
			require.Equal(t, exist, tc.exist)
		})
	}
}

func TestAccountRepository_ExistByEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewAccountRepository(pgTestSuite.db)

	testcases := []struct {
		name  string
		email string
		exist bool
	}{
		{
			name:  "exist",
			email: seed.AccountJohnDoe.Email,
			exist: true,
		},
		{
			name:  "wrong email",
			email: "wrong email",
			exist: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			exist, err := repo.ExistByEmail(ctx, tc.email)
			require.NoError(t, err)
			require.Equal(t, exist, tc.exist)
		})
	}
}
