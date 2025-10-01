package model

type SignUpModel struct {
	Username  string `form:"username" binding:"required"`
	Email     string `form:"email" binding:"required,email"`
	FirstName string `form:"first_name" binding:"required"`
	LastName  string `form:"last_name" binding:"required"`
	Password  string `form:"password" binding:"required"`
}
