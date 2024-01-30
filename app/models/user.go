package models

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
