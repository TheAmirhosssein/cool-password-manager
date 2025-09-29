package account

import "github.com/TheAmirhosssein/cool-password-manage/pkg/errors"

const (
	CodeAuthInvalidUser = 401_100
)

const (
	// Auth
	MessageAuthInvalidUser = "invalid user information"
)

var (
	// Auth
	AuthInvalidUser = errors.NewError(MessageAuthInvalidUser, CodeAuthInvalidUser)
)
