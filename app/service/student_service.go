package service

import (
	"PROJECTUAS_BE/app/repository"
	"PROJECTUAS_BE/middleware"

	"github.com/gofiber/fiber/v2"
)

type Studentservice struct {
	repo repository.StudentRepository
}

func NewAStudentService(repo repository.StudentRepository) *Studentservice {
	return &Studentservice{repo: repo}
}

func (s *Studentservice) GetStudent(c *fiber.Ctx) error {
	// Ambil id dari params
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	// Hanya admin yang bisa mengakses data student
	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)
	if userClaims.Role != "mahasiswa" {
		return fiber.NewError(fiber.StatusForbidden, "Access denied: mahasiswa only")
	}

	// Ambil data student berdasarkan user ID
	student, err := s.repo.GetStudentByUserID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Student not found")
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   student,
	})
}

func (s *Studentservice) GetAllStudents(c *fiber.Ctx) error {
	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)
	if userClaims.Role != "mahasiswa" {
		return fiber.NewError(fiber.StatusForbidden, "Access denied: mahasiswa only")
	}

	students, err := s.repo.GetAllStudents()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status":   "success",
		"students": students,
	})
}
