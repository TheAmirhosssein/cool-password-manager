package model

type SignUpInitModel struct {
	Username            string `json:"username" binding:"required"`
	Email               string `json:"email" binding:"required,email"`
	FirstName           string `json:"firstName" binding:"required"`
	LastName            string `json:"lastName" binding:"required"`
	RegistrationRequest []byte `json:"registrationRequest" binding:"required"`
}

type SignUpFinalizeModel struct {
	RegistrationID     string `form:"registrationID" binding:"required"`
	RegistrationRecord string `form:"registrationRecord" binding:"required"`
}

type LoginModel struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type LoginInitModel struct {
	Username string `json:"username" binding:"required"`
	KE1      []byte `json:"ke1" binding:"required"`
}

type TwoFactorModel struct {
	VerificationCode string `form:"verification_code" binding:"required"`
}
