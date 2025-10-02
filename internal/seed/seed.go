package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateSeed(ctx context.Context, db *pgxpool.Pool) {
	createAccountSeed(ctx, db)
	createGroupSeed(ctx, db)
}
