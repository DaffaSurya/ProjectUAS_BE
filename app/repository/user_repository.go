package repository

import (
	model "PROJECTUAS_BE/app/Model"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	FindByEmail(Email string) (*model.User, error)
	CreateUser(username, email, password, roleID, fullname string) (string, error)
	GetProfile(id string) (*model.User, error)
	GetAllUsers() ([]model.User, error)
	GetRoleByUserID(userID string) (string, error)
	GetRoleNameByRoleID(roleID string) (string, error)
	GetUserByID(id string) (*model.User, error)
	UpdateUserByID(id string, name string, email string) error
	DeleteUserByID(id string) error
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

func (r *userPostgres) CreateUser(username, email, password, roleID, fullname string) (string, error) {
	query := `
		INSERT INTO users (username, email, password_hash, role_id, full_name)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var userID string
	err := r.db.QueryRow(query, username, email, password, roleID, fullname).Scan(&userID)
	if err != nil {
		fmt.Println("DB ERROR:", err)
		return "", err
	}

	return userID, nil
}

func (r *userPostgres) GetRoleNameByRoleID(roleID string) (string, error) {
	query := "SELECT name FROM roles WHERE id = $1 LIMIT 1"

	var name string
	err := r.db.QueryRow(query, roleID).Scan(&name)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("role not found")
	}

	if err != nil {
		return "", err
	}

	return name, nil
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

func (r *userPostgres) GetRoleByUserID(userID string) (string, error) {
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

func (r *userPostgres) GetAllUsers() ([]model.User, error) {
	query := `
		SELECT id, username, email, password_hash, role_id, full_name, is_active
		FROM users;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		u := model.User{}
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.Password,
			&u.RoleID,
			&u.Fullname,
			&u.Is_active,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *userPostgres) GetUserByID(id string) (*model.User, error) {

	query := `
		SELECT  id, username, email, password_hash, role_id, full_name, is_active
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	user := model.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email,
		&user.Password,
		&user.RoleID,
		&user.Fullname,
		&user.Is_active)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userPostgres) UpdateUserByID(id string, name string, email string) error {
	query := `
        UPDATE users
        SET username = $1,
            email = $2,
            updated_at = NOW()
        WHERE id = $3
    `

	res, err := r.db.Exec(query, name, email, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userPostgres) DeleteUserByID(id string) error {
	query := `
        DELETE FROM users
        WHERE id = $1
    `

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
