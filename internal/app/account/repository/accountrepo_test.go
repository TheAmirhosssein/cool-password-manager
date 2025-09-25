package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/database"
	"github.com/TheAmirhosssein/cool-password-manage/internal/seed"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/testdocker"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

type postgresTest struct {
	db   *pgxpool.Pool
	name string
}

var pgTestSuite postgresTest

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgTask := make(chan postgresTest)
	go func() {
		pgName, pgTest := database.SetupTestDB(ctx)
		pgTask <- postgresTest{name: pgName, db: pgTest}
	}()

	pgTestSuite = <-pgTask

	seed.CreateSeed(ctx, pgTestSuite.db)

	exitCode := m.Run()
	testdocker.StopAndRemoveContainer(ctx, pgTestSuite.name, pgTestSuite.name)

	os.Exit(exitCode)
}

func TestAccountRepository_ExistByPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	testcases := []struct {
		name     string
		username string
		password string
		exist    bool
	}{
		{
			name:     "exist",
			username: seed.UserJohnDoe.Username,
			password: seed.UserJohnDoe.Password,
			exist:    true,
		},
		{
			name:     "wrong password",
			username: seed.UserJohnDoe.Username,
			password: "wrong password",
			exist:    false,
		},
		{
			name:     "wrong username",
			username: "wrong username",
			password: seed.UserJohnDoe.Password,
			exist:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Setup
			repo := repository.NewAccountRepository(pgTestSuite.db)

			// Read
			exist, err := repo.ExistByPassword(ctx, tc.username, tc.password)
			require.NoError(t, err)
			require.Equal(t, exist, tc.exist)
		})
	}
}
