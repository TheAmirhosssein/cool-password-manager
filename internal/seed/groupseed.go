package seed

import (
	"context"
	"fmt"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/base"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	idGroupBrockhampton types.ID = iota + 1
	idGroupOddFuture
	idGroupRadiohead
)

var (
	GroupBrockhampton = entity.Group{
		Entity:      base.Entity{ID: idGroupBrockhampton},
		Name:        "Brockhampton",
		Description: types.NullString{String: "Brockhampton Band Members", Valid: true},
		Owner:       AccountKevinAbstract,
		Members: []entity.Account{
			AccountMattChampion,
			AccountKevinAbstract,
			AccountJoba,
		},
	}

	GroupOddFuture = entity.Group{
		Entity:      base.Entity{ID: idGroupOddFuture},
		Name:        "Odd Future",
		Description: types.NullString{String: "Odd Future Band Members", Valid: true},
		Owner:       AccountTyler,
		Members: []entity.Account{
			AccountTyler,
			AccountEarl,
			AccountFrankOcean,
		},
	}

	GroupRadiohead = entity.Group{
		Entity:      base.Entity{ID: idGroupRadiohead},
		Name:        "Radiohead",
		Description: types.NullString{String: "Radiohead Band Members", Valid: true},
		Owner:       AccountThomYorke,
		Members: []entity.Account{
			AccountThomYorke,
			AccountJonnyGreenwood,
			AccountColinGreenwood,
		},
	}
)

func createGroupSeed(ctx context.Context, db *pgxpool.Pool) {
	var groups = []entity.Group{
		GroupBrockhampton,
		GroupOddFuture,
		GroupRadiohead,
	}

	for _, g := range groups {
		query := `
		INSERT INTO groups(name, description, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id;
		`
		var groupID int
		err := db.QueryRow(ctx, query, g.Name, g.Description.String, g.Owner.Entity.ID).Scan(&groupID)
		if err != nil {
			panic(fmt.Errorf("failed to insert group %s: %w", g.Name, err))
		}

		for _, m := range g.Members {
			_, err := db.Exec(ctx,
				"INSERT INTO groups_accounts(account_id, group_id) VALUES ($1, $2)",
				m.Entity.ID, groupID)
			if err != nil {
				panic(fmt.Errorf("failed to add member %s to group %s: %w", m.Username, g.Name, err))
			}
		}
	}
}
