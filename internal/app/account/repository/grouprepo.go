package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/param"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupRepository interface {
	Create(ctx context.Context, group *entity.Group) error
	Read(ctx context.Context, param param.ReadGroupParams) ([]entity.Group, int, error)
	ReadOne(ctx context.Context, id, memberID types.ID) (entity.Group, error)
	Update(ctx context.Context, group entity.Group) error
	AddAccounts(ctx context.Context, groupID types.ID, accounts []entity.Account) error
	DeleteAllMembers(ctx context.Context, groupID, ownerID types.ID) error
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

func (repo groupRepo) Read(ctx context.Context, param param.ReadGroupParams) ([]entity.Group, int, error) {
	query := `
	WITH paged_groups AS (
		SELECT g.id, g.name, g.description, g.owner_id
		FROM groups g
		WHERE g.id IN (
			SELECT group_id FROM groups_accounts WHERE account_id = $1
		)
		ORDER BY g.id
		LIMIT $2 OFFSET $3
	),
	rows_count AS (
		SELECT COUNT(*) AS count FROM groups g
		WHERE g.id IN (
			SELECT group_id FROM groups_accounts WHERE account_id = $1
		)
	)
	SELECT
    	rc.count, pg.id AS group_id, pg.name AS group_name, pg.description,
		o.id AS owner_id, o.username AS owner_username, o.first_name AS owner_first_name, 
		o.last_name AS owner_last_name, o.email AS owner_email,
		m.id as member_id, m.username AS member_username, m.first_name AS member_first_name,
		m.last_name AS member_last_name, m.email AS member_email
	FROM paged_groups pg
	JOIN accounts o ON o.id = pg.owner_id
	JOIN groups_accounts ga ON ga.group_id = pg.id
	JOIN accounts m ON m.id = ga.account_id
	CROSS JOIN rows_count rc
	ORDER BY group_id, member_id ASC;
	`

	rows, err := repo.db.Query(ctx, query, param.MemberID, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var count int
	groupMap := make(map[types.ID]*entity.Group)
	groupOrder := make([]types.ID, 0)

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
			groupOrder = append(groupOrder, groupID)
		} else {
			existing.Members = append(existing.Members, member)
		}
	}

	groups := make([]entity.Group, 0, len(groupOrder))
	for _, id := range groupOrder {
		groups = append(groups, *groupMap[id])
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
	WHERE g.id = $1 AND g.id IN (
		SELECT group_id FROM groups_accounts WHERE account_id = $2
	)
	ORDER BY g.id
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

func (repo groupRepo) Update(ctx context.Context, group entity.Group) error {
	query := "UPDATE groups SET name = $1, description = $2 WHERE id = $3 AND owner_id = $4"

	_, err := repo.db.Exec(ctx, query, group.Name, group.Description, group.Entity.ID, group.Owner.Entity.ID)
	if err != nil {
		log.ErrorLogger.Error("error at updating group", "error", err.Error())
		return err
	}

	return nil
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

func (repo groupRepo) DeleteAllMembers(ctx context.Context, groupID, ownerID types.ID) error {
	query := `
	DELETE FROM groups_accounts ga
	USING groups g
	WHERE g.id = ga.group_id AND g.id = $1 AND g.owner_id = $2
	`

	_, err := repo.db.Exec(ctx, query, groupID, ownerID)
	return err
}
