package repository

import (
	model "PROJECTUAS_BE/app/Model"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	FindByEmail(Email string) (*model.User, error)
	CreateUser(username, email, password, roleID, fullname string) error
	GetProfile(id string) (*model.User, error)
}

type userPostgres struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userPostgres{db}
}

func (r *userPostgres) FindByEmail(Email string) (*model.User, error) {
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

func (r *userPostgres) GetUserRole(userID int) (string, error) {
	var role string
	query := `
        SELECT roles.name
        FROM user_roles
        JOIN roles ON roles.id = user_roles.role_id
        WHERE user_roles.user_id=$1
    `
	err := r.db.QueryRow(query, userID).Scan(&role)
	return role, err
}

func (r *userPostgres) GetPermissionsByRole(roleName string) ([]string, error) {
	query := `
        SELECT p.name, p.resource, p.action
        FROM permissions p
        JOIN role_permissions rp ON rp.permission_id = p.id
        JOIN roles r ON r.id = rp.role_id
        WHERE r.name=$1
    `
	rows, err := r.db.Query(query, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var p string
		rows.Scan(&p)
		perms = append(perms, p)
	}

	return perms, nil
}

func (r *userPostgres) CreateUser(username, email, password, roleID, fullname string) error {
	query := `
		INSERT INTO users (username, email, password_hash, role_id, full_name)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, username, email, password, roleID, fullname)
	return err
}

func (r *userPostgres) GetProfile(id string) (*model.User, error) {
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
