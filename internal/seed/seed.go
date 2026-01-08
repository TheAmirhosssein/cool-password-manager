package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func CreateSeed(ctx context.Context, db *pgxpool.Pool) {
	createAccountSeed(ctx, db)
	createGroupSeed(ctx, db)
}

func CreateRedisSeed(ctx context.Context, redis *redis.Client) {
	createRegistrationSeed(ctx, redis)
	createRegistrationSeed(ctx, redis)
}
