package model

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password_hash"`
	Fullname string `json:"full_name"`
	RoleID   uint   `json: "role_id"` // foreign ke table roles
	Is_active bool `json: "is_active"`
}
