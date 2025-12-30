package usecase_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/usecase"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/database"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/totp"
	"github.com/TheAmirhosssein/cool-password-manage/internal/seed"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/encrypt"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/testdocker"
	"github.com/alicebob/miniredis/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	googleTotp "github.com/pquerna/otp/totp"
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

func TestAuthUsecase_SignUp(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	johnDoe := seed.AccountJohnDoe

	testcases := []struct {
		name        string
		account     entity.Account
		expectedErr error
	}{
		{
			name: "success signup",
			account: entity.Account{
				Username:  "new_user",
				Email:     "new_user@example.com",
				FirstName: "New",
				LastName:  "User",
			},
			expectedErr: nil,
		},
		{
			name: "username already exists",
			account: entity.Account{
				Username:  johnDoe.Username, // already seeded
				Email:     "unique_email@example.com",
				FirstName: "Dup",
				LastName:  "User",
			},
			expectedErr: account.AuthUsernameExist,
		},
		{
			name: "email already exists",
			account: entity.Account{
				Username:  "unique_username",
				Email:     johnDoe.Email, // already seeded
				FirstName: "Dup",
				LastName:  "User",
			},
			expectedErr: account.AuthEmailExist,
		},
		{
			name: "invalid password",
			account: entity.Account{
				Username:  "new_user",
				Email:     "new_user@example.com",
				FirstName: "New",
				LastName:  "User",
			},
			expectedErr: account.AuthInvalidPassword,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			u := setupAuthUsecase()
			auth, err := u.SignUp(ctx, tc.account)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				require.Equal(t, totp.Authenticator{}, auth)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, auth.Secret)
			}
		})
	}
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
			name:     "success",
			username: johnDoe.Username,
			// password:    seed.AccountDefaultPassword,
			expectedErr: nil,
			expectTF:    true,
		},
		{
			name:        "user does not exist",
			username:    "ghost",
			password:    "irrelevant",
			expectedErr: account.AuthInvalidAccount,
			expectTF:    false,
		},
		{
			name:        "invalid password",
			username:    johnDoe.Username,
			password:    "wrong password",
			expectedErr: account.AuthInvalidAccount,
			expectTF:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			u := setupAuthUsecase()
			tf, err := u.CreateTwoFactor(ctx, tc.username)

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

func TestAuthUsecase_ValidateTwoFactor(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	johnDoe := seed.AccountJohnDoe
	twoFactor := seed.TwoFactorJohnDoe // seeded in `seed.CreateTowFactorSeed`

	// setup usecase with real authenticator
	u := setupAuthUsecase()

	// decrypt JohnDoe’s secret so we can generate a valid TOTP code
	key, err := config.GetTestConfig().GetAESSecretKey()
	require.NoError(t, err)
	secret, err := encrypt.DecryptAESSecret(key, johnDoe.TOTPSecret)
	require.NoError(t, err)

	validCode, err := googleTotp.GenerateCode(secret, time.Now())
	require.NoError(t, err)

	testcases := []struct {
		name        string
		twoFactorID types.CacheID
		code        string
		expectedErr error
		expectAcc   bool
	}{
		{
			name:        "success login",
			twoFactorID: twoFactor.ID,
			code:        validCode,
			expectedErr: nil,
			expectAcc:   true,
		},
		{
			name:        "two factor does not exist",
			twoFactorID: "nonexistent-id",
			code:        "123456",
			expectedErr: account.AuthTwoFactorDoesNotExist,
			expectAcc:   false,
		},
		{
			name:        "invalid verification code",
			twoFactorID: twoFactor.ID,
			code:        "000000",
			expectedErr: account.AuthInvalidVerificationCode,
			expectAcc:   false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			acc, err := u.ValidateTwoFactor(ctx, tc.twoFactorID, tc.code)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				require.Equal(t, entity.Account{}, acc)
			} else {
				require.NoError(t, err)
				require.True(t, tc.expectAcc)
				require.Equal(t, johnDoe.Username, acc.Username)
				require.NotEmpty(t, acc.TOTPSecret)
			}
		})
	}
}

func setupAuthUsecase() usecase.AuthUsecase {
	aRepo := repository.NewAccountRepository(pgTestSuite.db)
	tfRepo := repository.NewTwoFactorRepository(redisClient)
	authenticator := totp.NewAuthenticatorAdaptor("something")
	config := config.GetTestConfig()

	return usecase.NewAuthUsecase(aRepo, tfRepo, authenticator, config)
}
