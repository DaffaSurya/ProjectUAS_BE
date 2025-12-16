package service

import (
	"PROJECTUAS_BE/app/repository"
	"PROJECTUAS_BE/middleware"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo repository.UserRepository
	repo repository.StudentPostgres
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) LoginService(email, password string) (string, error) {

	// Ambil user berdasarkan email
	user, err := s.Repo.FindByEmail(email)
	if err != nil || user == nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	// Cek password bcrypt
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	// Ambil role user berdasarkan tabel role_permissions
	role, err := s.Repo.GetRoleByUserID(user.ID)
	// role, err := s.Repo.GetRoleNameByRoleID(user.ID)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user role")
	}

	// Generate JWT termasuk role
	token, err := middleware.GenerateToken(
		user.ID,
		user.Username,
		user.Email,
		role,       // <-- pastikan diisi role
		[]string{}, // permissions jika diperlukan
	)

	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}

	return token, nil
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password_hash"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Panggil logic
	token, err := s.LoginService(body.Email, body.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token})
}

func (s *AuthService) GetProfile(c *fiber.Ctx) error {
	// Ambil claims dari middleware
	claimsData := c.Locals("claims")

	fmt.Println("DEBUG Authorization Header:", c.Get("Authorization")) // debugging session
	fmt.Println("DEBUG c.Locals(\"claims\"):", claimsData)

	if claimsData == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	claims, ok := claimsData.(*middleware.Claims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token data"})
	}

	// Query database by user_id
	user, err := s.Repo.GetProfile(claims.UserID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"username": user.Username,
			"fullname": user.Fullname,
			"email":    user.Email,
		},
	})
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	tokenStr := strings.Split(authHeader, " ")[1]

	claims := c.Locals("claims").(*middleware.Claims)

	// masukkan token ke blacklist dengan expiry JWT
	middleware.BlacklistToken(tokenStr, claims.ExpiresAt.Time)

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}

// ==================================
//
//	GET ALL USERS (ADMIN)
//
// ==================================
func (s *AuthService) GetAllUsers(c *fiber.Ctx) error {

	// Ambil claims dari middleware
	claims := c.Locals("claims")
	fmt.Println("CLAIMS IN SERVICE:", claims) // debugging
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	// Hanya admin yang bisa
	if userClaims.Role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Access denied: Admin only")
	}

	// Ambil semua user
	users, err := s.Repo.GetAllUsers()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}

func (s *AuthService) GetUsersByID(c *fiber.Ctx) error {

	// ambil id dari parameter repository
	id := c.Params("id")
	fmt.Println("UUID:", id)
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}
	// ambil id dari parameter repository

	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid uuid format")
	}

	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	// Hanya admin yang bisa
	if userClaims.Role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Access denied: Admin only")
	}

	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func (s *AuthService) CreateUser(c *fiber.Ctx) error {
	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	// Hanya admin yang bisa
	if userClaims.Role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Access denied: Admin only")
	}

	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password_hash"`
		RoleID   string `json:"role_id"`
		Fullname string `json:"full_name"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	// =============== Validasi ===============
	if body.Username == "" || body.Email == "" || body.Password == "" || body.RoleID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "All fields are required")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("error hash:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}

	userID, err := s.Repo.CreateUser(body.Username, body.Email, string(hashed), body.RoleID, body.Fullname)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}

	// 2️⃣ Get role name
	roleName, err := s.Repo.GetRoleNameByRoleID(body.RoleID)
	// roleName, err := s.Repo.GetRoleNameByRoleID(body.RoleID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	// 3️⃣ If student → insert into students table
	if roleName == "student" || roleName == "mahasiswa" {
		err := s.repo.CreateStudent(userID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "User created but failed to create student")
		}
	}

	return c.JSON(fiber.Map{
		"message": "User created successfully",
		"data": fiber.Map{
			"username": body.Username,
			"email":    body.Email,
			"role_id":  body.RoleID,
			"fullname": body.Fullname,
		},
	})

}

func (s *AuthService) UpdateUserByID(c *fiber.Ctx) error {
	// Ambil UUID dari URL
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	// Validasi format UUID
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid uuid format")
	}

	// Ambil claims dari middleware
	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	// Periksa apakah admin
	if userClaims.Role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Access denied: Admin only")
	}

	// Bind request body
	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Update melalui repository
	err := s.Repo.UpdateUserByID(id, body.Name, body.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "user updated successfully",
	})
}

func (s *AuthService) DeleteUserByID(c *fiber.Ctx) error {
	// Ambil UUID dari URL
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	// Validasi format UUID
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid uuid format")
	}

	// Ambil claims (dari JWT middleware)
	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	// Hanya admin yang boleh
	if userClaims.Role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Access denied: Admin only")
	}

	// Delete user melalui repository
	err := s.Repo.DeleteUserByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "user deleted successfully",
	})

}
