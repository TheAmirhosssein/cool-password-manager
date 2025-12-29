package entity

import (
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/base"
)

type Account struct {
	Entity     base.Entity
	Username   string
	Email      string
	FirstName  string
	LastName   string
	Password   string
	TOTPSecret []byte
}
