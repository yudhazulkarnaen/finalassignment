package responses

import "time"

type UserRegister struct {
	Age      uint   `json:"age" example:"23"`
	Email    string `json:"email" example:"name@org.dom.ge"`
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type UserUpdate struct {
	ID        uint      `json:"id" example:"1"`
	Email     string    `json:"email" example:"name@org.dom.ge"`
	Username  string    `json:"username"`
	Age       uint      `json:"age" example:"23"`
	UpdatedAt time.Time `json:"updated_at" example:"2019-11-09T21:21:46+00:00"`
}

type UserLogin struct {
	Token string `json:"token" example:"header.payload.signature"`
}
