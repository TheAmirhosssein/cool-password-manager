package http

const (
	// Me
	PathMe = "/"

	// Auth
	PathSignUp      = "/account/auth/sign-up/"
	PathSignUpInit  = "/account/auth/sign-up/init/"
	PathSignUpFinal = "/account/auth/sign-up/final/"
	PathLogin       = "/account/auth/login/"
	PathLoginInit   = "/account/auth/login/init/"
	PathTwoFactor   = "/account/auth/two-factor/"
	PathLogout      = "/account/auth/logout/"

	// Group
	PathGroupList         = "/account/groups/"
	PathGroupCreate       = "/account/groups/create/"
	PathGroupEdit         = "/account/groups/edit/"
	PathGroupDelete       = "/account/groups/delete/"
	PathGroupSearchMember = "/account/groups/members/"
)
