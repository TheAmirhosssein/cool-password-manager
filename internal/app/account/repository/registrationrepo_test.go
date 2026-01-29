package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/repository"
	"github.com/TheAmirhosssein/cool-password-manage/internal/seed"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/base"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRegistrationRepository_Create(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewRegistrationRepository(redisClient)

	testcases := []struct {
		name         string
		registration entity.Registration
		wantErr      bool
	}{
		{
			name: "successful",
			registration: entity.Registration{
				CacheEntity: base.CacheEntity{ID: "something", Duration: time.Minute},
				Username:    "something",
				Email:       "something@gmail.com",
				FirstName:   "something",
				LastName:    "something",
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := repo.Create(ctx, tc.registration)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				registration, err := repo.Get(ctx, tc.registration.ID)
				require.NoError(t, err)
				require.Equal(t, tc.registration.Username, registration.Username)
			}
		})
	}
}

func TestRegistrationRepository_Get(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := repository.NewRegistrationRepository(redisClient)

	tf := seed.RegistrationJohnDoe

	testcases := []struct {
		name     string
		id       types.CacheID
		expected entity.Registration
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

			registration, err := repo.Get(ctx, tc.id)
			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, registration.Username, tc.expected.Username)
				require.Equal(t, registration.Email, tc.expected.Email)
				require.Equal(t, registration.FirstName, tc.expected.FirstName)
				require.Equal(t, registration.LastName, tc.expected.LastName)
			}
		})
	}
}
