package repository

import (
	"context"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository interface {
	ReadByUsername(ctx context.Context, username string) (entity.Account, error)
	Update(ctx context.Context, account entity.Account) error
	ExistByPassword(ctx context.Context, username string, password string) (bool, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) AccountRepository {
	return repo{db: db}
}

func (r repo) ReadByUsername(ctx context.Context, username string) (entity.Account, error) {
	query := "SELECT username, email, secret FROM accounts WHERE username = $1"

	var account entity.Account
	err := r.db.QueryRow(ctx, query, username).Scan(&account.Username, &account.Email, &account.Secret)

	if err != nil {
		log.ErrorLogger.Error("getting account by username", "error", err.Error(), "username", username)
		return entity.Account{}, err
	}

	return account, nil
}

func (r repo) Update(ctx context.Context, account entity.Account) error {
	query := "UPDATE accounts SET secret = $1 WHERE id = $2"

	_, err := r.db.Exec(ctx, query, account.Secret, account.Entity.ID)
	if err != nil {
		log.ErrorLogger.Error("error at updating account", "error", err.Error(), "username", account.Username)
		return err
	}

	return nil
}

func (r repo) ExistByPassword(ctx context.Context, username string, password string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 WHERE username = $1 AND password = $2) FROM accounts"

	var exist bool
	err := r.db.QueryRow(ctx, query, username, password).Scan(&exist)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		log.ErrorLogger.Error("error at checking account existent by password and username", "error", err.Error(), "username", username)
		return false, err
	}

	return exist, nil
}
