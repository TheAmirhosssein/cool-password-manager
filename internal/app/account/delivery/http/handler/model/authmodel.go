package model

type SignUpModel struct {
	Username  string `form:"username" binding:"required"`
	Email     string `form:"email" binding:"required,email"`
	FirstName string `form:"first_name" binding:"required"`
	LastName  string `form:"last_name" binding:"required"`
	Password  string `form:"password" binding:"required"`
}

type LoginModel struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type TwoFactorModel struct {
	VerificationCode string `form:"verification_code" binding:"required"`
}
