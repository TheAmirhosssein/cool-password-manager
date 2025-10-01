package account

import "github.com/TheAmirhosssein/cool-password-manage/pkg/errors"

const (
	CodeAuthInvalidAccount = 401_100

	CodeAuthUsernameExist = 409_100
	CodeAuthEmailExist    = 409_101
)

const (
	// Auth
	MessageAuthInvalidAccount = "invalid user information"
	MessageAuthUsernameExist  = "taken username"
	MessageAuthEmailExist     = "an account with that email already exist"
)

var (
	// Auth
	AuthInvalidAccount = errors.NewError(MessageAuthInvalidAccount, CodeAuthInvalidAccount)
	AuthUsernameExist  = errors.NewError(MessageAuthUsernameExist, CodeAuthUsernameExist)
	AuthEmailExist     = errors.NewError(MessageAuthEmailExist, CodeAuthEmailExist)
)
