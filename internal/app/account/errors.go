package account

import "github.com/TheAmirhosssein/cool-password-manage/pkg/errors"

const (
	CodeAuthInvalidAccount = 401_100

	CodeAuthTwoFactorDoesNotExist = 404_100

	CodeAuthUsernameExist = 409_100
	CodeAuthEmailExist    = 409_101

	CodeAuthInvalidVerificationCode = 422_200
)

const (
	// Auth
	MessageAuthInvalidAccount          = "invalid user information"
	MessageAuthUsernameExist           = "taken username"
	MessageAuthEmailExist              = "an account with that email already exist"
	MessageAuthTwoFactorDoesNotExist   = "two factor authentication does not exist"
	MessageAuthInvalidVerificationCode = "the verification code is invalid"
)

var (
	// Auth
	AuthInvalidAccount          = errors.NewError(MessageAuthInvalidAccount, CodeAuthInvalidAccount)
	AuthUsernameExist           = errors.NewError(MessageAuthUsernameExist, CodeAuthUsernameExist)
	AuthEmailExist              = errors.NewError(MessageAuthEmailExist, CodeAuthEmailExist)
	AuthTwoFactorDoesNotExist   = errors.NewError(MessageAuthTwoFactorDoesNotExist, CodeAuthTwoFactorDoesNotExist)
	AuthInvalidVerificationCode = errors.NewError(MessageAuthInvalidVerificationCode, CodeAuthInvalidVerificationCode)
)
