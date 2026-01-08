package base

import (
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
)

type Entity struct {
	ID        types.ID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CacheEntity struct {
	ID       types.CacheID
	Duration time.Duration
}
