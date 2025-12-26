package account

import "github.com/TheAmirhosssein/cool-password-manage/pkg/errors"

const (
	CodeGroupInvalidGroupID = 400_100

	CodeAuthInvalidAccount = 401_100

	CodeGroupOnlyTheOwnerCanEdit   = 403_100
	CodeGroupOnlyTheOwnerCanDelete = 403_101

	CodeAuthTwoFactorDoesNotExist   = 404_100
	CodeAccountUsernameDoesNotExist = 404_101
	CodeGroupDoesNotExist           = 404_102

	CodeAuthUsernameExist = 409_100
	CodeAuthEmailExist    = 409_101

	CodeAuthInvalidPassword         = 422_100
	CodeAuthInvalidVerificationCode = 422_101
)

const (
	// Auth
	MessageAuthInvalidAccount          = "invalid user information"
	MessageAuthUsernameExist           = "taken username"
	MessageInvalidPassword             = "invalid password"
	MessageAuthEmailExist              = "an account with that email already exist"
	MessageAuthTwoFactorDoesNotExist   = "two factor authentication does not exist"
	MessageAuthInvalidVerificationCode = "the verification code is invalid"

	// Group
	MessageGroupOnlyTheOwnerCanEdit   = "only the group owner can edit the group"
	MessageGroupOnlyTheOwnerCanDelete = "only the group owner can delete the group"
	MessageGroupInvalidGroupID        = "invalid group id"
	MessageGroupDoesNotExist          = "group does not exist"

	// Account
	MessageAccountUsernameDoesNotExist = "account with that username does not exist"
)

var (
	// Auth
	AuthInvalidAccount          = errors.NewError(MessageAuthInvalidAccount, CodeAuthInvalidAccount)
	AuthUsernameExist           = errors.NewError(MessageAuthUsernameExist, CodeAuthUsernameExist)
	AuthEmailExist              = errors.NewError(MessageAuthEmailExist, CodeAuthEmailExist)
	AuthInvalidPassword         = errors.NewError(MessageInvalidPassword, CodeAuthInvalidPassword)
	AuthTwoFactorDoesNotExist   = errors.NewError(MessageAuthTwoFactorDoesNotExist, CodeAuthTwoFactorDoesNotExist)
	AuthInvalidVerificationCode = errors.NewError(MessageAuthInvalidVerificationCode, CodeAuthInvalidVerificationCode)

	// Group
	GroupOnlyTheOwnerCanEdit   = errors.NewError(MessageGroupOnlyTheOwnerCanEdit, CodeGroupOnlyTheOwnerCanEdit)
	GroupOnlyTheOwnerCanDelete = errors.NewError(MessageGroupOnlyTheOwnerCanDelete, CodeGroupOnlyTheOwnerCanDelete)
	GroupInvalidGroupID        = errors.NewError(MessageGroupInvalidGroupID, CodeGroupInvalidGroupID)
	GroupDoesNotExist          = errors.NewError(MessageGroupDoesNotExist, CodeGroupDoesNotExist)

	// Account
	AccountUsernameDoesNotExist = errors.NewError(MessageAccountUsernameDoesNotExist, CodeAccountUsernameDoesNotExist)
)
