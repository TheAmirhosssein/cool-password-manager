package seed

import (
	"context"
	"encoding/json"
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/base"
	"github.com/redis/go-redis/v9"
)

const (
	idJohnDoeRegistration = "john_doe_registration"
)

const (
	DefaultPassword = "strong-password"
)

var (
	RegistrationJohnDoe = entity.Registration{
		CacheEntity: base.CacheEntity{ID: idJohnDoeRegistration, Duration: time.Minute},
		Username:    AccountJohnDoe.Username,
		Email:       AccountJohnDoe.Email,
		FirstName:   AccountJohnDoe.FirstName,
		LastName:    AccountJohnDoe.LastName,
	}
)

func createRegistrationSeed(ctx context.Context, rdb *redis.Client) {
	marshaledRegistration, err := json.Marshal(RegistrationJohnDoe)
	if err != nil {
		panic(err)
	}

	err = rdb.Set(ctx, string(RegistrationJohnDoe.ID), marshaledRegistration, RegistrationJohnDoe.Duration).Err()
	if err != nil {
		panic(err)
	}
}
