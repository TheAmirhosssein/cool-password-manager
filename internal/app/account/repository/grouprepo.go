package repository

import (
	"context"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupRepository interface {
	Create(ctx context.Context, group entity.Group) error
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
