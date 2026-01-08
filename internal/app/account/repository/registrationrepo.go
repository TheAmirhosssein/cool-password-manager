package repository

import (
	"context"
	"encoding/json"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/redis/go-redis/v9"
)

type RegistrationRepository interface {
	Create(ctx context.Context, registration entity.Registration) error
	Get(ctx context.Context, id types.CacheID) (entity.Registration, error)
}

type registrationRepo struct {
	client *redis.Client
}

func NewRegistrationRepository(client *redis.Client) RegistrationRepository {
	return registrationRepo{client: client}
}

func (r registrationRepo) Create(ctx context.Context, registration entity.Registration) error {
	marshaledRegistration, err := json.Marshal(registration)
	if err != nil {
		log.ErrorLogger.Error("error marshaling registration", "error", err.Error())
		return err
	}

	err = r.client.Set(ctx, string(registration.ID), marshaledRegistration, registration.Duration).Err()
	if err != nil {
		log.ErrorLogger.Error("error saving registration", "error", err.Error(), "username", registration.Username)
		return err
	}

	return nil
}

func (r registrationRepo) Get(ctx context.Context, id types.CacheID) (entity.Registration, error) {
	result, err := r.client.Get(ctx, string(id)).Bytes()
	if err != nil {
		log.ErrorLogger.Error("error getting registration", "error", err.Error(), "id", id)
		return entity.Registration{}, err
	}

	registration := new(entity.Registration)
	if err := json.Unmarshal(result, registration); err != nil {
		log.ErrorLogger.Error("error at unmarshaling registration", "error", err.Error())
		return entity.Registration{}, err
	}

	return *registration, nil
}
