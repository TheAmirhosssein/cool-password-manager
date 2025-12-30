package seed

import (
	"context"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/base"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	idAccountJohnDoe types.ID = iota + 1
	idUserMattChampion
	idUserKevinAbstract
	idUserJoba
	idUserTyler
	idUserEarl
	idUserFrankOcean
	idUserKendrickLamar
	idUserJayRock
	idUserColinGreenwood
)

const (
	AccountDefaultPassword = "123"
)

var (
	AccountJohnDoe = entity.Account{
		Entity:     base.Entity{ID: idAccountJohnDoe},
		Username:   "j.doe",
		Email:      "j.doe@gmail.com",
		FirstName:  "John",
		LastName:   "Doe",
		TOTPSecret: []byte("UDkdflLm0Z6yaRIKJnEAb3dndEVPRsdIx3V6CmKJ49ihhoybL8m157tPyAs7l6Cm8rfyME50UHr9dxbE"),
	}

	AccountMattChampion = entity.Account{
		Entity:    base.Entity{ID: idUserMattChampion},
		Username:  "m.champion",
		Email:     "m.champion@gmail.com",
		FirstName: "Matt",
		LastName:  "Champion",
	}

	AccountKevinAbstract = entity.Account{
		Entity:    base.Entity{ID: idUserKevinAbstract},
		Username:  "k.abstract",
		Email:     "k.abstract@gmail.com",
		FirstName: "Kevin",
		LastName:  "Abstract",
	}

	AccountJoba = entity.Account{
		Entity:    base.Entity{ID: idUserJoba},
		Username:  "d.joba",
		Email:     "d.joba@gmail.com",
		FirstName: "Joba (Russel)",
		LastName:  "Boring",
	}

	AccountTyler = entity.Account{
		Entity:    base.Entity{ID: idUserTyler},
		Username:  "tyler",
		Email:     "tyler@gmail.com",
		FirstName: "Tyler",
		LastName:  "The Creator",
	}

	AccountEarl = entity.Account{
		Entity:    base.Entity{ID: idUserEarl},
		Username:  "earl",
		Email:     "earl@gmail.com",
		FirstName: "Earl",
		LastName:  "Sweatshirt",
	}

	AccountFrankOcean = entity.Account{
		Entity:    base.Entity{ID: idUserFrankOcean},
		Username:  "frank",
		Email:     "frank@gmail.com",
		FirstName: "Frank",
		LastName:  "Ocean",
	}

	AccountKendrickLamar = entity.Account{
		Entity:    base.Entity{ID: idUserKendrickLamar},
		Username:  "k.lamar",
		Email:     "k.lamar@gmail.com",
		FirstName: "Kendrick",
		LastName:  "Lamar",
	}

	AccountJayRock = entity.Account{
		Entity:    base.Entity{ID: idUserJayRock},
		Username:  "j.rock",
		Email:     "j.rock@gmail.com",
		FirstName: "Jay",
		LastName:  "Rock",
	}

	AccountSchoolBoyQ = entity.Account{
		Entity:    base.Entity{ID: idUserColinGreenwood},
		Username:  "schoolboy.q",
		Email:     "schoolboy.q@gmail.com",
		FirstName: "SchoolBoy",
		LastName:  "Q",
	}

	AccountAbSoul = entity.Account{
		Entity:    base.Entity{ID: idUserColinGreenwood},
		Username:  "a.soul",
		Email:     "a.soulq@gmail.com",
		FirstName: "ab",
		LastName:  "soul",
	}
)

func createAccountSeed(ctx context.Context, db *pgxpool.Pool) {
	query := `
	INSERT INTO accounts(username, email, first_name, last_name, totp_secret) VALUES
	('j.doe', 'j.doe@gmail.com', 'John', 'Doe', 'UDkdflLm0Z6yaRIKJnEAb3dndEVPRsdIx3V6CmKJ49ihhoybL8m157tPyAs7l6Cm8rfyME50UHr9dxbE'),
	('m.champion', 'm.champion@gmail.com', 'Matt', 'Champion',  ''),
	('k.abstract', 'k.abstract@gmail.com', 'Kevin', 'Abstract',  ''),
	('d.joba', 'd.joba@gmail.com', 'Dom', 'Joba',  ''),
	('tyler', 'tyler@gmail.com', 'Tyler', 'The Creator',  ''),
	('earl', 'earl@gmail.com', 'Earl', 'Sweatshirt',  ''),
	('frank', 'frank@gmail.com', 'Frank', 'Ocean',  ''),
	('k.lamar', 'k.lamar@gmail.com', 'Kendrick', 'Lamar',  ''),
	('j.rock', 'j.rock@gmail.com', 'Jay', 'Rock',  ''),
	('schoolboy.q', 'schoolboy.q@gmail.com', 'SchoolBoy', 'Q',  ''),
	('a.soul', 'a.soul@gmail.com', 'Ab', 'Soul',  '');
	`
	_, err := db.Exec(ctx, query)
	if err != nil {
		panic(err)
	}
}
