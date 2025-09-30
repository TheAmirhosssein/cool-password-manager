package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/seed"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestTwoFactorRepository_Set(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewTwoFactorRepository(redisClient)

	testcases := []struct {
		name      string
		twoFactor entity.TwoFactor
		wantErr   bool
	}{
		{
			name:      "successful",
			twoFactor: entity.TwoFactor{ID: "something", Username: "Something", Duration: time.Minute},
			wantErr:   false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := repo.Create(ctx, tc.twoFactor)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				twoFactor, err := repo.Get(ctx, tc.twoFactor.ID)
				require.NoError(t, err)
				require.Equal(t, tc.twoFactor.Username, twoFactor.Username)
			}
		})
	}
}

func TestTwoFactorRepository_Get(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewTwoFactorRepository(redisClient)

	tf := seed.TwoFactorJohnDoe

	testcases := []struct {
		name     string
		id       types.CacheID
		expected entity.TwoFactor
		err      error
	}{
		{
			name:     "successful",
			id:       tf.ID,
			expected: tf,
		}, {
			name: "not found",
			err:  redis.Nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			twoFactor, err := repo.Get(ctx, tc.id)
			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, twoFactor.Username, tc.expected.Username)
			}
		})
	}
}

func TestTwoFactorRepository_Delete(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewTwoFactorRepository(redisClient)

	testcases := []struct {
		name     string
		id       types.CacheID
		expected entity.TwoFactor
		err      error
	}{
		{
			name: "successful",
			id:   seed.TwoFactorJohnDoe.ID,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := repo.Delete(ctx, tc.id)
			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				exist, err := repo.Exist(ctx, tc.id)
				require.NoError(t, err)
				require.False(t, exist)
			}
		})
	}
}

func TestTwoFactorRepository_Exist(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewTwoFactorRepository(redisClient)

	testcases := []struct {
		name  string
		id    types.CacheID
		exist bool
		err   error
	}{
		{
			name:  "successful",
			id:    seed.TwoFactorJohnDoe.ID,
			exist: true,
		},
		{
			name:  "does not exist",
			id:    "wrong_id",
			exist: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			exist, err := repo.Exist(ctx, tc.id)
			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, exist, tc.exist)
			}
		})
	}
}
