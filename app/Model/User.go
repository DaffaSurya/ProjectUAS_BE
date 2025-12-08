package model

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password_hash"`
	RoleID    string `json: "role_id"`
	Fullname  string `json:"full_name"` // foreign ke table roles
	Is_active bool   `json:"is_active"`
}

type Profile struct {
	username string `json:"username"`
	fullname string `json:"full_name"`
	Email    string `json:"email"`
}
