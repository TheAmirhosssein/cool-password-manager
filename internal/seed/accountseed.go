package seed

import (
	"context"

	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/base"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	idUserJohnDoe types.ID = iota + 1
	idUserMattChampion
)

var (
	UserJohnDoe = entity.Account{
		Entity:    base.Entity{ID: idUserJohnDoe},
		Username:  "j.doe",
		Email:     "j.doe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "$2a$14$htOj2XrIL7y7bF8OrEbyneK9Nd20mHx7vYQ4fK8x/ZR6MJS0hQjWy", // 123
	}
	UserMattChampion = entity.Account{
		Entity:    base.Entity{ID: idUserMattChampion},
		Username:  "m.champion",
		Email:     "m.champion@gmail.com",
		FirstName: "Matt",
		LastName:  "Champion",
		Password:  "$2a$14$htOj2XrIL7y7bF8OrEbyneK9Nd20mHx7vYQ4fK8x/ZR6MJS0hQjWy", // 123
	}
)

func createAccountSeed(ctx context.Context, db *pgxpool.Pool) {
	query := `
	INSERT INTO accounts(username, email, first_name, last_name, password) 
	VALUES ('j.doe', 'j.doe@gmail.com', 'John', 'Doe', '$2a$14$htOj2XrIL7y7bF8OrEbyneK9Nd20mHx7vYQ4fK8x/ZR6MJS0hQjWy'),
		   ('m.champion', 'm.champion@gmail.com', 'Matt', 'Champion', '$2a$14$htOj2XrIL7y7bF8OrEbyneK9Nd20mHx7vYQ4fK8x/ZR6MJS0hQjWy')
	`
	_, err := db.Exec(ctx, query)
	if err != nil {
		panic(err)
	}
}
