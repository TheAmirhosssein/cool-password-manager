package entity

import "github.com/TheAmirhosssein/cool-password-manage/internal/utils/base"

type Registration struct {
	base.CacheEntity
	CredID    []byte `json:"cred_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
