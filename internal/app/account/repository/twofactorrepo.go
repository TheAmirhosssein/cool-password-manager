package repository

import (
	"context"
	"errors"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/redis/go-redis/v9"
)

type TwoFactorRepository interface {
	Create(ctx context.Context, twoFactor entity.TwoFactor) error
	Get(ctx context.Context, id types.CacheID) (entity.TwoFactor, error)
	Delete(ctx context.Context, id types.CacheID) error
	Exist(ctx context.Context, id types.CacheID) (bool, error)
}

type twoFactorRepo struct {
	client *redis.Client
}

func NewTwoFactorRepository(client *redis.Client) TwoFactorRepository {
	return twoFactorRepo{client: client}
}

func (r twoFactorRepo) Create(ctx context.Context, twoFactor entity.TwoFactor) error {
	err := r.client.Set(ctx, string(twoFactor.ID), twoFactor.Username, twoFactor.Duration).Err()
	if err != nil {
		log.ErrorLogger.Error("error saving two factor", "error", err.Error(), "username", twoFactor.Username)
		return err
	}

	return nil
}

func (r twoFactorRepo) Get(ctx context.Context, id types.CacheID) (entity.TwoFactor, error) {
	result, err := r.client.Get(ctx, string(id)).Result()
	if err != nil {
		log.ErrorLogger.Error("error getting two factor", "error", err.Error(), "id", id)
		return entity.TwoFactor{}, err
	}

	return entity.TwoFactor{ID: id, Username: result}, nil
}

func (r twoFactorRepo) Delete(ctx context.Context, id types.CacheID) error {
	err := r.client.Del(ctx, string(id)).Err()
	if err := err; err != nil {
		log.ErrorLogger.Error("error deleting two factor", "error", err.Error(), "id", id)
		return err
	}

	return err
}

func (r twoFactorRepo) Exist(ctx context.Context, id types.CacheID) (bool, error) {
	value, err := r.client.Get(ctx, string(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) && value == "" {
			return false, nil
		}
		log.ErrorLogger.Error("error checking two factor existence", "error", err.Error(), "id", id)
		return false, err
	}

	return true, nil
}
