package seed

import (
	"context"
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/redis/go-redis/v9"
)

const (
	idJohnDoeTwoFactor = "john_doe_two_factor"
)

var (
	TwoFactorJohnDoe = entity.TwoFactor{
		ID:       idJohnDoeTwoFactor,
		Username: AccountJohnDoe.Username,
		Duration: time.Minute * 2,
	}
)

func createTowFactorSeed(ctx context.Context, rdb *redis.Client) {
	err := rdb.Set(ctx, string(TwoFactorJohnDoe.ID), TwoFactorJohnDoe.Username, TwoFactorJohnDoe.Duration).Err()
	if err != nil {
		panic(err)
	}
}
