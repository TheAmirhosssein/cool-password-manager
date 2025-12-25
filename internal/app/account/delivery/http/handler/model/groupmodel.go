package model

import "github.com/TheAmirhosssein/cool-password-manage/internal/types"

type GroupUpdate struct {
	Name        string     `form:"name" binding:"required"`
	Description string     `form:"description" binding:"required"`
	MembersID   []types.ID `form:"members[]" binding:"required"`
}
