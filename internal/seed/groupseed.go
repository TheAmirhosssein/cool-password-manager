package seed

import (
	"context"

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
			AccountMattChampion,
		},
	}

	GroupBlackHippy = entity.Group{
		Entity:      base.Entity{ID: idGroupRadiohead},
		Name:        "Black Hippy",
		Description: types.NullString{String: "Black Hippy Band Members", Valid: true},
		Owner:       AccountKendrickLamar,
		Members: []entity.Account{
			AccountKendrickLamar,
			AccountJayRock,
			AccountSchoolBoyQ,
			AccountAbSoul,
		},
	}
)

func createGroupSeed(ctx context.Context, db *pgxpool.Pool) {
	gQuery := `
	INSERT INTO groups(name, description, owner_id)
	VALUES ('Brockhampton', 'Brockhampton Band Members', 3),
		('Odd Future', 'Odd Future Band Members', 5),
		('Black Hippy', 'Black Hippy', 8);
	`

	_, err := db.Exec(ctx, gQuery)
	if err != nil {
		panic(err)
	}

	gaQuery := `
	INSERT INTO groups_accounts(group_id, account_id)
	VALUES (1, 2), (1, 3), (1, 4), -- Brockhampton
		(2, 5), (2, 6), (2, 7), -- Odd Future
		(3, 8), (3, 9), (3, 10), (3, 11); -- Black Hippy

	`

	_, err = db.Exec(ctx, gaQuery)
	if err != nil {
		panic(err)
	}
}
