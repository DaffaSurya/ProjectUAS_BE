package service

import (
	"PROJECTUAS_BE/app/repository"
	"PROJECTUAS_BE/middleware"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo repository.UserRepository
	repo repository.StudentPostgres
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// ==================================
//
//	GET ALL USERS (ADMIN)
//
// ==================================
func (s *UserService) GetAllUsers(c *fiber.Ctx) error {

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

func (s *UserService) GetUsersByID(c *fiber.Ctx) error {

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

func (s *UserService) CreateUser(c *fiber.Ctx) error {
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

func (s *UserService) UpdateUserByID(c *fiber.Ctx) error {
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

func (s *UserService) DeleteUserByID(c *fiber.Ctx) error {
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
