package repository

import (
	model "PROJECTUAS_BE/app/Model"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
    user := &model.User{}

    query := `SELECT id, username, email, password_hash FROM users WHERE email=$1`

    err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (r *UserRepository) GetUserRole(userID int) (string, error) {
    var role string
    query := `
        SELECT roles.name
        FROM user_roles
        JOIN roles ON roles.id = user_roles.role_id
        WHERE user_roles.user_id=$1
    `
    err := r.DB.QueryRow(query, userID).Scan(&role)
    return role, err
}


func (r *UserRepository) GetPermissionsByRole(roleName string) ([]string, error) {
    query := `
        SELECT p.name, p.resource, p.action
        FROM permissions p
        JOIN role_permissions rp ON rp.permission_id = p.id
        JOIN roles r ON r.id = rp.role_id
        WHERE r.name=$1
    `
    rows, err := r.DB.Query(query, roleName)
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
