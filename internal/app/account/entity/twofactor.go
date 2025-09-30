package entity

import (
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
)

type TwoFactor struct {
	ID       types.CacheID
	Username string
	Duration time.Duration
}
