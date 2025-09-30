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
)

const (
	AccountJohnDoePassword = "123"
)

var (
	AccountJohnDoe = entity.Account{
		Entity:    base.Entity{ID: idAccountJohnDoe},
		Username:  "j.doe",
		Email:     "j.doe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou", // 123
	}
	AccountMattChampion = entity.Account{
		Entity:    base.Entity{ID: idUserMattChampion},
		Username:  "m.champion",
		Email:     "m.champion@gmail.com",
		FirstName: "Matt",
		LastName:  "Champion",
		Password:  "$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou", // 123
	}
)

func createAccountSeed(ctx context.Context, db *pgxpool.Pool) {
	query := `
	INSERT INTO accounts(username, email, first_name, last_name, password) 
	VALUES ('j.doe', 'j.doe@gmail.com', 'John', 'Doe', '$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou'),
		   ('m.champion', 'm.champion@gmail.com', 'Matt', 'Champion', '$2a$14$ygA10iotMn5KQQ46qQTpIOCFIzPawSyWuQ8oeh2FEUlFrbkqOiSou')
	`
	_, err := db.Exec(ctx, query)
	if err != nil {
		panic(err)
	}
}
