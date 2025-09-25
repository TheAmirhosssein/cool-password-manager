package repository

import (
	"context"

	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository interface {
	ExistByPassword(ctx context.Context, username string, password string) (bool, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) AccountRepository {
	return repo{db: db}
}

func (r repo) ExistByPassword(ctx context.Context, username string, password string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 WHERE username = $1 AND password = $2) FROM accounts"

	var exist bool
	err := r.db.QueryRow(ctx, query, username, password).Scan(&exist)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		log.ErrorLogger.Error("getting account by password and username failed", "error", err.Error(), "username", username)
		return false, err
	}

	return exist, nil
}
