package usecase_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
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

func TestAuthUsecase_CreateTwoFactor(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	johnDoe := seed.AccountJohnDoe

	testcases := []struct {
		name        string
		username    string
		password    string
		expectedErr error
		expectTF    bool
	}{
		{
			name:        "success",
			username:    johnDoe.Username,
			password:    seed.AccountJohnDoePassword,
			expectedErr: nil,
			expectTF:    true,
		},
		{
			name:        "user does not exist",
			username:    "ghost",
			password:    "irrelevant",
			expectedErr: account.AuthInvalidUser,
			expectTF:    false,
		},
		{
			name:        "invalid password",
			username:    johnDoe.Username,
			password:    "wrong password",
			expectedErr: account.AuthInvalidUser,
			expectTF:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			u := setupAuthUsecase()
			tf, err := u.CreateTwoFactor(ctx, tc.username, tc.password)

			// Assertions
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				require.Equal(t, entity.TwoFactor{}, tf)
			} else {
				require.NoError(t, err)
				require.True(t, tc.expectTF)
				require.Equal(t, tc.username, tf.Username)
				require.NotEmpty(t, tf.ID)
			}
		})
	}
}

func setupAuthUsecase() usecase.AuthUsecase {
	aRepo := repository.NewAccountRepository(pgTestSuite.db)
	tfRepo := repository.NewTwoFactorRepository(redisClient)
	return usecase.NewAuthUsecase(aRepo, tfRepo)
}
