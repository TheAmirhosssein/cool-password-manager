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
	ExistByUsername(ctx context.Context, username string) (bool, error)
}

type accountRepo struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) AccountRepository {
	return accountRepo{db: db}
}

func (r accountRepo) ReadByUsername(ctx context.Context, username string) (entity.Account, error) {
	query := "SELECT username, email, secret, password FROM accounts WHERE username = $1"

	var account entity.Account
	err := r.db.QueryRow(ctx, query, username).Scan(&account.Username, &account.Email, &account.Secret, &account.Password)

	if err != nil {
		log.ErrorLogger.Error("getting account by username", "error", err.Error(), "username", username)
		return entity.Account{}, err
	}

	return account, nil
}

func (r accountRepo) Update(ctx context.Context, account entity.Account) error {
	query := "UPDATE accounts SET secret = $1 WHERE id = $2"

	_, err := r.db.Exec(ctx, query, account.Secret, account.Entity.ID)
	if err != nil {
		log.ErrorLogger.Error("error at updating account", "error", err.Error(), "username", account.Username)
		return err
	}

	return nil
}

func (r accountRepo) ExistByUsername(ctx context.Context, username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 WHERE username = $1) FROM accounts"

	var exist bool
	err := r.db.QueryRow(ctx, query, username).Scan(&exist)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		log.ErrorLogger.Error("error at checking account existent by username", "error", err.Error(), "username", username)
		return false, err
	}

	return exist, nil
}
