package service

import (
	model "PROJECTUAS_BE/app/Model"
	"PROJECTUAS_BE/app/repository"
	"PROJECTUAS_BE/middleware"
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LecturesService struct {
	Repo repository.LecturesRepository
}

func NewLecturesService(repo repository.LecturesRepository) *LecturesService {
	return &LecturesService{Repo: repo}
}

func (s *LecturesService) VerifyAchievement(c *fiber.Ctx) error {
	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)
	if userClaims.Role != "dosen" {
		return fiber.NewError(
			fiber.StatusForbidden,
			"Only dosen wali can verify achievements",
		)
	}

	achievementID := c.Params("id")
	if achievementID == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"Achievement ID is required",
		)
	}

	var req model.VerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"Invalid request body",
		)
	}

	if req.Status != "verified" && req.Status != "rejected" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"Status must be 'verified' or 'rejected'",
		)
	}

	if req.Status == "rejected" && req.RejectionReason == nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"Rejection reason is required",
		)
	}

	err := s.Repo.Verify(
		context.Background(),
		achievementID,
		req.Status,
		userClaims.UserID,
		req.RejectionReason,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(
				fiber.StatusNotFound,
				"Achievement not found or not in submitted status",
			)
		}
		return fiber.NewError(
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}

	// ===== 6. Response =====
	return c.JSON(fiber.Map{
		"message": "Achievement verification successful",
		"data": fiber.Map{
			"mongo_achievement_id": achievementID,
			"status":               req.Status,
			"verified_at":          time.Now(),
			"verified_by":          userClaims.UserID,
			"rejection_reason":     req.RejectionReason,
		},
	})

}

func (s *LecturesService) RejectAchievement(c *fiber.Ctx) error {
	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)
	if userClaims.Role != "dosen" {
		return fiber.NewError(
			fiber.StatusForbidden,
			"Only dosen wali can verify achievements",
		)
	}

	achievementID := c.Params("id")
	if achievementID == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"Achievement ID is required",
		)
	}

	var req model.RejectRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"Invalid request body",
		)
	}

	err := s.Repo.Reject(
		context.Background(),
		achievementID,
		req.Reason,
		userClaims.UserID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(
				fiber.StatusNotFound,
				"Achievement not found or not in submitted status",
			)
		}
		return fiber.NewError(
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}

	// ===== 6. Response =====
	return c.JSON(fiber.Map{
		"message": "Achievement rejected successfully",
		"data": fiber.Map{
			"mongo_achievement_id": achievementID,
			"status":               "rejected",
			"reason":               req.Reason,
		},
	})
}

func (s *LecturesService) GetHistory(c *fiber.Ctx) error { // Get History bisa diakses oleh admin , student dan dosen wali

	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	// ===== 2. Param =====
	mongoAchievementID := c.Params("id")
	if mongoAchievementID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Achievement ID is required")
	}

	histories, err := s.Repo.GetHistory(
		context.Background(),
		mongoAchievementID,
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch achievement history")
	}

	if len(histories) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "No history found")
	}

	if userClaims.Role == "mahasiswa" {
		if histories[0].StudentID != userClaims.UserID {
			return fiber.NewError(fiber.StatusForbidden, "Access denied")
		}
	}
	return c.JSON(fiber.Map{ // response
		"message": "Achievement history fetched successfully",
		"data":    histories,
	})
}

func (s *LecturesService) GetLectures(c *fiber.Ctx) error {
	lecturers, err := s.Repo.GetallLectures(context.Background())
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"Failed to fetch lecturers",
		)
	}

	return c.JSON(fiber.Map{
		"message": "Lecturers fetched successfully",
		"total":   len(lecturers),
		"data":    lecturers,
	})
}

func (s *LecturesService) Getadvisees(c *fiber.Ctx) error {
	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	lecturerID := c.Params("id")
	if lecturerID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Lecturer ID is required")
	}

	switch userClaims.Role {
	case "admin", "mahasiswa":
		// allowed
	case "dosen":
		// dosen hanya boleh lihat mahasiswa bimbingannya sendiri
		if userClaims.UserID != lecturerID {
			return fiber.NewError(fiber.StatusForbidden, "Access denied")
		}
	default:
		return fiber.NewError(fiber.StatusForbidden, "Access denied")
	}

	// ===== 4. Get Data =====
	students, err := s.Repo.Getadvisees(
		context.Background(),
		lecturerID,
	)
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"Failed to get advisees",
		)
	}

	// ===== 5. Response =====
	return c.JSON(fiber.Map{
		"lecturer_id": lecturerID,
		"total":       len(students),
		"students":    students,
	})
}
