package model

import "github.com/TheAmirhosssein/cool-password-manage/internal/types"

type GroupUpdate struct {
	Name        string     `form:"name" binding:"required"`
	Description string     `form:"description"`
	MembersID   []types.ID `form:"members[]" binding:"required"`
}

type GroupCreate struct {
	Name        string     `form:"name" binding:"required"`
	Description string     `form:"description"`
	MembersID   []types.ID `form:"members[]" binding:"required"`
}
