package entity

import (
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/base"
)

type Group struct {
	base.Entity
	Name        string
	Description types.NullString
	Owner       Account
	Members     []Account
}
