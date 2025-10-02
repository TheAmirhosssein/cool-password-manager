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
	idUserThomYorke
	idUserJonnyGreenwood
	idUserColinGreenwood
)

const (
	AccountDefaultPassword = "123"
)

var (
	AccountJohnDoe = entity.Account{
		Entity:    base.Entity{ID: idAccountJohnDoe},
		Username:  "j.doe",
		Email:     "j.doe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou",
		Secret:    "UDkdflLm0Z6yaRIKJnEAb3dndEVPRsdIx3V6CmKJ49ihhoybL8m157tPyAs7l6Cm8rfyME50UHr9dxbE",
	}

	AccountMattChampion = entity.Account{
		Entity:    base.Entity{ID: idUserMattChampion},
		Username:  "m.champion",
		Email:     "m.champion@gmail.com",
		FirstName: "Matt",
		LastName:  "Champion",
		Password:  "$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou",
	}

	AccountKevinAbstract = entity.Account{
		Entity:    base.Entity{ID: idUserKevinAbstract},
		Username:  "k.abstract",
		Email:     "k.abstract@gmail.com",
		FirstName: "Kevin",
		LastName:  "Abstract",
		Password:  "$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou",
	}

	AccountJoba = entity.Account{
		Entity:    base.Entity{ID: idUserJoba},
		Username:  "d.joba",
		Email:     "d.joba@gmail.com",
		FirstName: "Joba (Russel)",
		LastName:  "Boring",
		Password:  "$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou",
	}

	AccountTyler = entity.Account{
		Entity:    base.Entity{ID: idUserTyler},
		Username:  "tyler",
		Email:     "tyler@gmail.com",
		FirstName: "Tyler",
		LastName:  "The Creator",
		Password:  "$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou",
	}

	AccountEarl = entity.Account{
		Entity:    base.Entity{ID: idUserEarl},
		Username:  "earl",
		Email:     "earl@gmail.com",
		FirstName: "Earl",
		LastName:  "Sweatshirt",
		Password:  "$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou",
	}

	AccountFrankOcean = entity.Account{
		Entity:    base.Entity{ID: idUserFrankOcean},
		Username:  "frank",
		Email:     "frank@gmail.com",
		FirstName: "Frank",
		LastName:  "Ocean",
		Password:  "$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou",
	}

	AccountThomYorke = entity.Account{
		Entity:    base.Entity{ID: idUserThomYorke},
		Username:  "thom",
		Email:     "thom@gmail.com",
		FirstName: "Thom",
		LastName:  "Yorke",
		Password:  "$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou",
	}

	AccountJonnyGreenwood = entity.Account{
		Entity:    base.Entity{ID: idUserJonnyGreenwood},
		Username:  "jonny",
		Email:     "jonny@gmail.com",
		FirstName: "Jonny",
		LastName:  "Greenwood",
		Password:  "$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou",
	}

	AccountColinGreenwood = entity.Account{
		Entity:    base.Entity{ID: idUserColinGreenwood},
		Username:  "colin",
		Email:     "colin@gmail.com",
		FirstName: "Colin",
		LastName:  "Greenwood",
		Password:  "$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou",
	}
)

func createAccountSeed(ctx context.Context, db *pgxpool.Pool) {
	query := `
	INSERT INTO accounts(username, email, first_name, last_name, password, secret) VALUES
	('j.doe', 'j.doe@gmail.com', 'John', 'Doe', '$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou', 'UDkdflLm0Z6yaRIKJnEAb3dndEVPRsdIx3V6CmKJ49ihhoybL8m157tPyAs7l6Cm8rfyME50UHr9dxbE'),
	('m.champion', 'm.champion@gmail.com', 'Matt', 'Champion', '$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou', ''),
	('k.abstract', 'k.abstract@gmail.com', 'Kevin', 'Abstract', '$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou', ''),
	('d.joba', 'd.joba@gmail.com', 'Dom', 'Joba', '$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou', ''),
	('tyler', 'tyler@gmail.com', 'Tyler', 'The Creator', '$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou', ''),
	('earl', 'earl@gmail.com', 'Earl', 'Sweatshirt', '$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou', ''),
	('frank', 'frank@gmail.com', 'Frank', 'Ocean', '$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou', ''),
	('thom', 'thom@gmail.com', 'Thom', 'Yorke', '$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou', ''),
	('jonny', 'jonny@gmail.com', 'Jonny', 'Greenwood', '$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou', ''),
	('colin', 'colin@gmail.com', 'Colin', 'Greenwood', '$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou', '')
	`
	_, err := db.Exec(ctx, query)
	if err != nil {
		panic(err)
	}
}
