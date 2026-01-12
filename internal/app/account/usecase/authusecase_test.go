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
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/opaque"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/totp"
	"github.com/TheAmirhosssein/cool-password-manage/internal/seed"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/encrypt"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/errors"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/testdocker"
	"github.com/alicebob/miniredis/v2"
	bytemareOpaque "github.com/bytemare/opaque"
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
var conf *config.Config

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

	conf = config.GetTestConfig()

	seed.CreateRedisSeed(ctx, redisClient)

	exitCode := m.Run()
	testdocker.StopAndRemoveContainer(ctx, pgTestSuite.name, pgTestSuite.name)

	os.Exit(exitCode)
}

func TestAuthUsecase_SignUpInit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	johnDoe := seed.AccountJohnDoe

	// valid opaque registration init message
	conf := bytemareOpaque.DefaultConfiguration()
	client, err := bytemareOpaque.NewClient(conf)
	require.NoError(t, err)

	password := []byte("strong-password")

	message := client.RegistrationInit(password).Serialize()

	testcases := []struct {
		name        string
		reg         entity.Registration
		message     []byte
		expectedErr error
	}{
		{
			name: "success sign up init",
			reg: entity.Registration{
				Username: "new_user",
				Email:    "new_user@example.com",
			},
			message:     message,
			expectedErr: nil,
		},
		{
			name: "username already exists",
			reg: entity.Registration{
				Username: johnDoe.Username,
				Email:    "unique_email@example.com",
			},
			message:     message,
			expectedErr: account.AuthUsernameExist,
		},
		{
			name: "email already exists",
			reg: entity.Registration{
				Username: "unique_username",
				Email:    johnDoe.Email,
			},
			message:     message,
			expectedErr: account.AuthEmailExist,
		},
		{
			name: "invalid opaque message",
			reg: entity.Registration{
				Username: "another_user",
				Email:    "another@example.com",
			},
			message:     []byte("invalid-message"),
			expectedErr: errors.NewServerError(),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			u := setupAuthUsecase()

			resp, registrationID, err := u.SignUpInit(ctx, tc.reg, tc.message)

			if tc.expectedErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErr.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, resp)

				// verify registration saved in redis
				registrationRepo := repository.NewRegistrationRepository(redisClient)
				saved, err := registrationRepo.Get(ctx, registrationID)
				require.NoError(t, err)
				require.Equal(t, tc.reg.Username, saved.Username)
				require.NotEmpty(t, saved.CredID)
				require.NotZero(t, saved.Email)
			}
		})
	}
}

func TestAuthUsecase_SignUpFinalize(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	u := setupAuthUsecase()

	// ---------- OPAQUE client setup ----------
	opaqueConf := bytemareOpaque.DefaultConfiguration()
	client, err := bytemareOpaque.NewClient(opaqueConf)
	require.NoError(t, err)

	password := []byte("strong-password")

	// ---------- SignUpInit phase ----------
	initMsg := client.RegistrationInit(password).Serialize()

	reg := entity.Registration{
		Username:  "new_user",
		Email:     "new_user@example.com",
		FirstName: "New",
		LastName:  "User",
	}

	resp, registrationID, err := u.SignUpInit(ctx, reg, initMsg)
	require.NoError(t, err)
	require.NotEmpty(t, resp)

	response, err := client.Deserialize.RegistrationResponse(resp)
	if err != nil {
		log.Fatalln(err)
	}

	record, _ := client.RegistrationFinalize(response, bytemareOpaque.ClientRegistrationFinalizeOptions{
		ClientIdentity: []byte(reg.Username),
		ServerIdentity: []byte(conf.Opaque.ServerID),
	})
	message3 := record.Serialize()

	testcases := []struct {
		name           string
		message        []byte
		registrationID types.CacheID
		expectedErr    bool
	}{
		{
			name:           "success signup finalize",
			message:        message3,
			registrationID: registrationID,
			expectedErr:    false,
		},
		{
			name:           "registration does not exist",
			message:        message3,
			registrationID: "non-existent-id",
			expectedErr:    true,
		},
		{
			name:           "invalid opaque message",
			message:        []byte("invalid-message"),
			registrationID: registrationID,
			expectedErr:    true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			auth, err := u.SignUpFinalize(ctx, tc.message, tc.registrationID)

			if tc.expectedErr {
				require.Error(t, err)
				require.Equal(t, totp.Authenticator{}, auth)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, auth.Secret)
			require.NotEmpty(t, auth.QrCode)

			// ---------- Verify account persisted ----------
			accRepo := repository.NewAccountRepository(pgTestSuite.db)
			acc, err := accRepo.ReadByUsername(ctx, reg.Username)
			require.NoError(t, err)

			require.Equal(t, reg.Username, acc.Username)
			require.Equal(t, reg.Email, acc.Email)
			require.NotEmpty(t, acc.OpaqueRecord)
			require.NotEmpty(t, acc.TOTPSecret)
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
	rRepo := repository.NewRegistrationRepository(redisClient)
	authenticator := totp.NewAuthenticatorAdaptor("something")
	opqaue, err := opaque.New(conf)
	if err != nil {
		panic(err)
	}

	return usecase.NewAuthUsecase(aRepo, tfRepo, rRepo, authenticator, opqaue, conf)
}
