package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	params "github.com/TheAmirhosssein/cool-password-manage/internal/app/param"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupRepository interface {
	Create(ctx context.Context, group *entity.Group) error
	Read(ctx context.Context, param params.ReadGroupParams) ([]entity.Group, int64, error)
	ReadOne(ctx context.Context, id, memberID types.ID) (entity.Group, error)
	AddAccounts(ctx context.Context, groupID types.ID, accounts []entity.Account) error
	RemoveAccounts(ctx context.Context, groupID types.ID, accounts []entity.Account) error
}

type groupRepo struct {
	db *pgxpool.Pool
}

func NewGroupRepository(db *pgxpool.Pool) GroupRepository {
	return groupRepo{db: db}
}

func (repo groupRepo) Create(ctx context.Context, group *entity.Group) error {
	query := "INSERT INTO groups (name, description, owner_id) VALUES ($1, $2, $3) RETURNING id"

	err := repo.db.QueryRow(ctx, query, group.Name, group.Description, group.Owner.Entity.ID).Scan(&group.ID)
	if err != nil {
		log.ErrorLogger.Error("error at creating group", "error", err.Error())
		return err
	}

	return nil
}

func (repo groupRepo) Read(ctx context.Context, param params.ReadGroupParams) ([]entity.Group, int64, error) {
	query := `
	WITH data AS (
		SELECT g.id, g.name, g.description,
				o.id, o.username, o.first_name, o.last_name, o.email,
				m.id, m.username, m.first_name, m.last_name, m.email
		FROM groups g
		JOIN accounts o ON o.id = g.owner_id
		JOIN groups_accounts ga ON ga.group_id = g.id
		JOIN accounts m ON m.id = ga.account_id
		WHERE g.id IN (
			SELECT group_id FROM groups_accounts WHERE account_id = $1
		)
		ORDER BY g.id
		LIMIT $2 OFFSET $3
	),
	rows_count AS (
		SELECT COUNT(g.id) AS count
		FROM groups g
		WHERE g.id IN (
			SELECT group_id FROM groups_accounts WHERE account_id = $1
		)
	)
	SELECT rows_count.count, data.* FROM data CROSS JOIN rows_count
	`

	rows, err := repo.db.Query(ctx, query, param.MemberID, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var count int64
	groupMap := make(map[types.ID]*entity.Group)

	for rows.Next() {
		var (
			groupID types.ID
			g       entity.Group
			owner   entity.Account
			member  entity.Account
		)

		err := rows.Scan(
			&count, &groupID, &g.Name, &g.Description,
			&owner.Entity.ID, &owner.Username, &owner.FirstName, &owner.LastName, &owner.Email,
			&member.Entity.ID, &member.Username, &member.FirstName, &member.LastName, &member.Email,
		)
		if err != nil {
			return nil, 0, err
		}

		existing, ok := groupMap[groupID]
		if !ok {
			g.Entity.ID = groupID
			g.Owner = owner
			g.Members = []entity.Account{member}
			groupMap[groupID] = &g
		} else {
			existing.Members = append(existing.Members, member)
		}
	}

	groups := make([]entity.Group, 0, len(groupMap))
	for _, g := range groupMap {
		groups = append(groups, *g)
	}

	return groups, count, nil
}

func (repo groupRepo) ReadOne(ctx context.Context, id, memberID types.ID) (entity.Group, error) {
	query := `
	SELECT g.id, g.name, g.description,
	       o.id, o.username, o.first_name, o.last_name, o.email,
	       m.id, m.username, m.first_name, m.last_name, m.email
	FROM groups g
	JOIN accounts o ON o.id = g.owner_id
	JOIN groups_accounts ga ON ga.group_id = g.id
	JOIN accounts m ON m.id = ga.account_id
	WHERE g.id = $1 AND ga.account_id = $2
	`

	rows, err := repo.db.Query(ctx, query, id, memberID)
	if err != nil {
		return entity.Group{}, err
	}

	var g entity.Group
	for rows.Next() {
		var member entity.Account
		err := rows.Scan(
			&g.Entity.ID, &g.Name, &g.Description,
			&g.Owner.Entity.ID, &g.Owner.Username, &g.Owner.FirstName, &g.Owner.LastName, &g.Owner.Email,
			&member.Entity.ID, &member.Username, &member.FirstName, &member.LastName, &member.Email,
		)
		if err != nil {
			return entity.Group{}, err
		}

		g.Members = append(g.Members, member)
	}

	return g, nil
}

func (repo groupRepo) AddAccounts(ctx context.Context, groupID types.ID, accounts []entity.Account) error {
	if len(accounts) == 0 {
		return nil
	}

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
	if err != nil {
		log.ErrorLogger.Error("error at creating group members", "error", err.Error())
		return err
	}

	return nil
}

func (repo groupRepo) RemoveAccounts(ctx context.Context, groupID types.ID, accounts []entity.Account) error {
	var (
		args   []any
		params []string
	)

	args = append(args, groupID)

	for i, account := range accounts {
		args = append(args, account.Entity.ID)
		params = append(params, fmt.Sprintf("$%d", i+2))
	}

	query := fmt.Sprintf(
		`DELETE FROM groups_accounts 
		 WHERE group_id = $1 AND account_id IN (%s)`,
		strings.Join(params, ","),
	)

	_, err := repo.db.Exec(ctx, query, args...)
	return err
}
