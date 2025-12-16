package service

import (
	model "PROJECTUAS_BE/app/Model"
	"PROJECTUAS_BE/app/repository"
	"PROJECTUAS_BE/middleware"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (s *Studentservice) SubmitAchievement(c *fiber.Ctx) error {
	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)
	if userClaims.Role != "mahasiswa" {
		return fiber.NewError(fiber.StatusForbidden, "Access denied: mahasiswa only")
	}

	fmt.Println("Id user:", userClaims.UserID)

	achievementID := c.Params("id")
	if achievementID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Achievement ID is required")
	}

	studentid, err := s.repo.GetStudentIDByUserID(
		context.Background(),
		userClaims.UserID,
	)
	if err != nil {
		return fiber.NewError(
			fiber.StatusNotFound,
			"Student not found",
		)
	}

	// ===== 3. Prepare Data =====
	now := time.Now()

	ref := &model.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          studentid,
		MongoAchievementID: achievementID,
		Status:             "submitted",
		SubmittedAt:        &now,
	}

	err = s.repo.Submit(context.Background(), ref)
	if err != nil {
		log.Println("ERROR SUBMIT ACHIEVEMENT:", err)
		return fiber.NewError(
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Achievement submitted successfully",
		"data": fiber.Map{
			"id":                   ref.ID,
			"mongo_achievement_id": ref.MongoAchievementID,
			"status":               ref.Status,
			"submitted_at":         ref.SubmittedAt,
		},
	})
}
