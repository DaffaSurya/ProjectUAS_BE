package service

import (
	"PROJECTUAS_BE/app/repository"
	"PROJECTUAS_BE/middleware"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) RegisterService(username, email, password, roleID, fullname string) error {

	// Cek apakah username sudah dipakai
	_, err := s.Repo.FindByEmail(username)
	if err == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Username already used")
	}

	// Hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Insert user
	return s.Repo.CreateUser(username, email, string(hashed), roleID, fullname)
}

func (s *AuthService) Register(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password_hash"`
		RoleID   string `json:"role_id"`
		Fullname string `json:"full_name"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := s.RegisterService(body.Username, body.Email, body.Password, body.RoleID, body.Fullname)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Registration successful"})
}

func (s *AuthService) LoginService(Email, password string) (string, error) {

	user, err := s.Repo.FindByEmail(Email)
	if err != nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid email")
	}

	// Cek password hash
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		fmt.Println("bcrypt error:", err)
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Generate JWT
	token, err := middleware.GenerateToken(
		user.ID,
		user.Username,
		user.Email,
		user.RoleID,
		[]string{}, // future: permissions
	)

	return token, err
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
