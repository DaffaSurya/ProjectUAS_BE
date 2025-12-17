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
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) LoginService(email, password string) (string, error) {

	// Ambil user berdasarkan email
	user, err := s.repo.FindByEmail(email)
	if err != nil || user == nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	// Cek password bcrypt
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	// Ambil role user berdasarkan tabel role_permissions
	role, err := s.repo.GetRoleByUserID(user.ID)
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
	user, err := s.repo.GetProfile(claims.UserID)
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
