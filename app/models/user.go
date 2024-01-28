package models

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
	GuestRole Role = "guest"
)

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Role     Role   `json:"role" form:"role"`
}
