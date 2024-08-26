package auth

type LoginRequestValidate struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
	AppId    int    `validate:"required"`
}
