package service

import (
	"PROJECTUAS_BE/app/repository"
	"PROJECTUAS_BE/middleware"
	"context"

	"github.com/gofiber/fiber/v2"
)

type ReportService struct {
	Repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) *ReportService {
	return &ReportService{
		Repo: repo,
	}
}

func (s *ReportService) GetStatics(c *fiber.Ctx) error {
	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	filter := repository.StatisticsFilter{}

	switch userClaims.Role {
	case "mahasiswa":
		filter.StudentID = &userClaims.UserID

	case "lecturer":
		filter.LecturerID = &userClaims.UserID

	case "admin":
		// no filter

	default:
		return fiber.NewError(fiber.StatusForbidden, "Access denied")
	}

	stats, err := s.Repo.GetStatics(
		context.Background(),
		filter,
	)
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"Failed to fetch statistics",
		)
	}

	return c.JSON(stats)
}

func (s *ReportService) GetStudentReport(c *fiber.Ctx) error {
	claims := c.Locals("claims")
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	studentID := c.Params("id")

	if studentID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Student ID required")
	}

	switch userClaims.Role {

	case "mahasiswa":
		if userClaims.UserID != studentID {
			return fiber.NewError(fiber.StatusForbidden, "Access denied")
		}

	case "lecturer":
		isAdvisor, err := s.Repo.IsAdvisor(
			context.Background(),
			userClaims.UserID,
			studentID,
		)
		if err != nil || !isAdvisor {
			return fiber.NewError(fiber.StatusForbidden, "Access denied")
		}

	case "admin":
		// allow

	default:
		return fiber.NewError(fiber.StatusForbidden, "Access denied")
	}

	report, err := s.Repo.GetStudentReport(
		context.Background(),
		studentID,
	)
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"Failed to get student report",
		)
	}

	return c.JSON(report)

}
