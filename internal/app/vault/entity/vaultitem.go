package entity

import (
	"github.com/TheAmirhosssein/cool-password-manage/internal/app/account/entity"
	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/base"
)

type ValueItem struct {
	base.Entity
	Name              string
	Description       types.NullString
	EncryptedUsername byte
	EncryptedPassword byte
	EncryptedUrl      types.NullByte
	EncryptedNote     types.NullByte
	Nonce             byte
	Creator           entity.Account
	Groups            []entity.Group
}
