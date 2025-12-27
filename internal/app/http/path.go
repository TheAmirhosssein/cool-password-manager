package http

const (
	// Me
	PathMe = "/"

	// Auth
	PathSignUp    = "/account/auth/sign-up/"
	PathLogin     = "/account/auth/login/"
	PathTwoFactor = "/account/auth/two-factor/"
	PathLogout    = "/account/auth/logout/"

	// Group
	PathGroupList         = "/account/groups/"
	PathGroupCreate       = "/account/groups/create/"
	PathGroupEdit         = "/account/groups/edit/"
	PathGroupDelete       = "/account/groups/delete/"
	PathGroupSearchMember = "/account/groups/members/"
)
