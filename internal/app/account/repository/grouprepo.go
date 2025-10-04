package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupRepository interface {
	Create(ctx context.Context, group entity.Group) error
	AddAccount(ctx context.Context, groupID types.ID, accounts []entity.Account) error
}

type groupRepo struct {
	db *pgxpool.Pool
}

func NewGroupRepository(db *pgxpool.Pool) GroupRepository {
	return groupRepo{db: db}
}

func (repo groupRepo) Create(ctx context.Context, group entity.Group) error {
	query := "INSERT INTO groups (name, description, owner_id) VALUES ($1, $2, $3)"

	_, err := repo.db.Exec(ctx, query, group.Name, group.Description, group.Owner.Entity.ID)
	if err != nil {
		log.ErrorLogger.Error("error at creating group", "error", err.Error())
		return err
	}

	return nil
}

func (repo groupRepo) AddAccount(ctx context.Context, groupID types.ID, accounts []entity.Account) error {
	var (
		values []string
		args   []any
	)

	args = append(args, groupID)

	for i, account := range accounts {
		placeholder := fmt.Sprintf("($1, $%d)", i+2)
		values = append(values, placeholder)

		args = append(args, account.Entity.ID)
	}

	query := fmt.Sprintf(
		"INSERT INTO groups_accounts (group_id, account_id) VALUES %s",
		strings.Join(values, ","),
	)

	_, err := repo.db.Exec(ctx, query, args...)
	return err
}
