package seed

import (
	"context"
	"encoding/base64"

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
	DefaultPassword = "strong-password"
)

var (
	AccountJohnDoe = entity.Account{
		Entity:       base.Entity{ID: idAccountJohnDoe},
		Username:     "j.doe",
		Email:        "j.doe@gmail.com",
		FirstName:    "John",
		LastName:     "Doe",
		TOTPSecret:   []byte("UDkdflLm0Z6yaRIKJnEAb3dndEVPRsdIx3V6CmKJ49ihhoybL8m157tPyAs7l6Cm8rfyME50UHr9dxbE"),
		OpaqueRecord: []byte("AsOOBVMFNKcXOaeZUL6ty1ybQL5IArnwI9tBuQxPiWqOfevEHtH4gKXCyb/Rc6ThXatGVqvzwnMmJskmFX27S+yLTcNP/4LiCtF+9muK6sXSZ9/Xx1z8URXn9ib39EB+eUBgA+kRTVVZ4e+wl5h8poZsn+c529/gwmea1LlSZYZ7"),
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
	record, err := base64.StdEncoding.DecodeString(string(AccountJohnDoe.OpaqueRecord))
	if err != nil {
		panic(err)
	}

	query := `
	INSERT INTO accounts
	(username, email, first_name, last_name, opaque_record, totp_secret)
	VALUES
	('j.doe', 'j.doe@gmail.com', 'John', 'Doe', $1,
	 'UDkdflLm0Z6yaRIKJnEAb3dndEVPRsdIx3V6CmKJ49ihhoybL8m157tPyAs7l6Cm8rfyME50UHr9dxbE'),

	('m.champion', 'm.champion@gmail.com', 'Matt', 'Champion',
	 'M0rjZ9F1x1F0YxRjM6Y1ZKq5A2V+8vD+JY4H7xX2V9k=',
	 ''),

	('k.abstract', 'k.abstract@gmail.com', 'Kevin', 'Abstract',
	 'f4rVQmJc6p8E3D0xK8K0M4Q5E1Zz9XJ+5B2K6p4n2zY=',
	 ''),

	('d.joba', 'd.joba@gmail.com', 'Dom', 'Joba',
	 'pFZ2X9H5W4m8r2XJZP9QKZc1X8T4r0mF8ZP9cW1xVY=',
	 ''),

	('tyler', 'tyler@gmail.com', 'Tyler', 'The Creator',
	 'ZK1P9K5xXQ2mY5N1X3Z9ZpJ5D8W0H2xY5c2T1Z9P0A=',
	 ''),

	('earl', 'earl@gmail.com', 'Earl', 'Sweatshirt',
	 'J5X9D8H0KZP1Y5Z2N3Q8P9W0Z5X1K2mF5cR4T1YV0A=',
	 ''),

	('frank', 'frank@gmail.com', 'Frank', 'Ocean',
	 'X9PZ5Z1J8K2N0H5D4Y5mR1X0Q3Z5W2cP9F8VY1A=',
	 ''),

	('k.lamar', 'k.lamar@gmail.com', 'Kendrick', 'Lamar',
	 'Z5X1Q2mF8R9P0H5D4Y5ZP9W1K2N0cX8J5VY1A=',
	 ''),

	('j.rock', 'j.rock@gmail.com', 'Jay', 'Rock',
	 'P0H5D4Y5Z5X1Q2mF8R9W1K2N0cX8J5VY1A=',
	 ''),

	('schoolboy.q', 'schoolboy.q@gmail.com', 'SchoolBoy', 'Q',
	 'X5ZP9W1K2N0H5D4Y5Z5X1Q2mF8R9cX8J5VY1A=',
	 ''),

	('a.soul', 'a.soul@gmail.com', 'Ab', 'Soul',
	 'Q2mF8R9P0H5D4Y5Z5X1K2N0ZP9W1cX8J5VY1A=',
	 '');
	`
	_, err = db.Exec(ctx, query, record)
	if err != nil {
		panic(err)
	}
}
