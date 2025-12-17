package repository

import (
	model "PROJECTUAS_BE/app/Model"
	"database/sql"
	"fmt"
)

type AuthRepository interface {
	FindByEmail(Email string) (*model.User, error)
	GetRoleByUserID(userID string) (string, error)
	GetProfile(id string) (*model.User, error)
}

type AuthPostGres struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &AuthPostGres{db: db}
}

func (r *AuthPostGres) FindByEmail(Email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, role_id, full_name
		FROM users
		WHERE email = $1
		LIMIT 1;
	`

	row := r.db.QueryRow(query, Email)

	user := model.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.Fullname,
	)

	fmt.Println("FindByEmail error:", err)
	fmt.Println("User result:", user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthPostGres) GetRoleByUserID(userID string) (string, error) {
	query := `
        SELECT ro.name
        FROM users u
        JOIN roles ro ON u.role_id = ro.id
        WHERE u.id = $1
        LIMIT 1;
    `

	var role string
	err := r.db.QueryRow(query, userID).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
}

func (r *AuthPostGres) GetProfile(id string) (*model.User, error) {
	query := `
		SELECT id, username, email, full_name
		FROM users
		WHERE id = $1
		LIMIT 1;
	`

	user := new(model.User)

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Fullname,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
