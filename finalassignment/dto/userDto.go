package dto

type UserRegister struct {
	Username string `validate:"required"`
	Email    string `validate:"required,email" example:"name@org.dom.ge"`
	Password string `validate:"required,min=6"`
	Age      uint   `validate:"required,gt=8" example:"23"`
}
type UserLogin struct {
	Email    string `validate:"required,email" example:"name@org.dom.ge"`
	Password string `validate:"required,min=6"`
}
type UserUpdate struct {
	Username string `validate:"required" json:"username"`
	Email    string `validate:"required,email" json:"email" example:"name@org.dom.ge"`
}
