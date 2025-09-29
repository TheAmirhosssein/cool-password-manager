package entity

import "github.com/TheAmirhosssein/cool-password-manage/internal/types"

type TwoFactor struct {
	ID       types.CacheID
	Username string
}
