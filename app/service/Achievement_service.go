package service

import (
	model "PROJECTUAS_BE/app/Model"
	"PROJECTUAS_BE/app/repository"
	"PROJECTUAS_BE/middleware"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type AchievementService struct {
	Repo repository.AchievementRepository
}

func NewAchievementService(repo repository.AchievementRepository) *AchievementService {
	return &AchievementService{
		Repo: repo,
	}
}

func (s *AchievementService) GetAllAchievements(c *fiber.Ctx) error {

	// Panggil repository langsung (service disatukan disini)
	achievements, err := s.Repo.GetAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get achievements")
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   achievements,
	})
}

func (s *AchievementService) GetAchievementsByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid achievement id")
	}

	// Panggil repository langsung
	achievement, err := s.Repo.GetAchievementByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch achievement")
	}

	if achievement == nil {
		return fiber.NewError(fiber.StatusNotFound, "achievement not found")
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   achievement,
	})
}

func (s *AchievementService) CreateAchievements(c *fiber.Ctx) error {
	// Ambil role dari JWT middleware
	claims := c.Locals("claims")
	studentId := c.Locals("user_id")
	fmt.Println("CLAIMS IN SERVICE:", claims) // debugging
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	// Hanya admin yang bisa
	if userClaims.Role != "mahasiswa" {
		return fiber.NewError(fiber.StatusForbidden, "Access ditolak, hanya mahasiswa yang boleh mengakses")
	}

	// Bind input
	input := new(model.Achievement)
	if err := c.BodyParser(input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Set required fields
	input.ID = uuid.New().String()
	input.StudentID = studentId.(string)
	input.CreatedAt = time.Now()

	err := s.Repo.Create(context.Background(), input)

	if err != nil {
		fmt.Println("ERROR SAVE ACHIEVEMENT:", err) // debug log

		return fiber.NewError(fiber.StatusInternalServerError,
			fmt.Sprintf("Failed to save achievement: %v", err))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Achievement added successfully",
		"data":    input,
	})
}

func (s *AchievementService) UpdateAchievement(c *fiber.Ctx) error {

	// AUTH: hanya mahasiswa
	// role := c.Locals("role")
	claims := c.Locals("claims")
	userId := c.Locals("user_id")

	fmt.Println("CLAIMS IN SERVICE:", claims) // debugging
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	// Hanya admin yang bisa
	if userClaims.Role != "mahasiswa" {
		return fiber.NewError(fiber.StatusForbidden, "Access ditolak, hanya mahasiswa yang boleh mengakses")
	}

	// GET ID FROM PARAM
	id := c.Params("id")

	// Cek apakah data ada
	existing, err := s.Repo.FindById(context.Background(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Achievement not found")
	}

	// AUTHORIZATION: hanya pemilik data yg boleh edit
	if existing.StudentID != userId {
		return fiber.NewError(fiber.StatusForbidden, "You cannot edit someone else's achievement")
	}

	// Bind request body
	req := new(model.Achievement)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Data yang boleh diupdate
	update := bson.M{
		"title":           req.Title,
		"Tags":            req.Tags,
		"AchievementType": req.AchievementType,
		"point":           req.Points,
		"description":     req.Description,
		"updated_at":      time.Now(),
	}

	// Lakukan update
	err = s.Repo.Update(context.Background(), id, update)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update achievement")
	}

	return c.JSON(fiber.Map{
		"message": "Achievement updated successfully",
		"id":      id,
	})
}

func (s *AchievementService) DeleteAchievement(c *fiber.Ctx) error {

	claims := c.Locals("claims")

	fmt.Println("CLAIMS IN SERVICE:", claims) // debugging
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	// Hanya admin yang bisa
	if userClaims.Role != "mahasiswa" {
		return fiber.NewError(fiber.StatusForbidden, "Access ditolak, hanya mahasiswa yang boleh mengakses")
	}

	id := c.Params("id")

	// Cek apakah data exist
	achievement, err := s.Repo.GetAchievementByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "database error")
	}

	if achievement == nil {
		return fiber.NewError(fiber.StatusNotFound, "achievement not found")
	}

	// Hapus
	err = s.Repo.Delete(context.Background(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete achievement")
	}

	return c.JSON(fiber.Map{
		"message": "achievement deleted successfully",
	})
}

func (s *AchievementService) GetStudentAchievements(c *fiber.Ctx) error {

	claims := c.Locals("claims")

	fmt.Println("CLAIMS IN SERVICE:", claims) // debugging
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	// Hanya admin yang bisa
	if userClaims.Role != "mahasiswa" {
		return fiber.NewError(fiber.StatusForbidden, "Access ditolak, hanya mahasiswa yang boleh mengakses")
	}

	studentID := userClaims.UserID
	if studentID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Student ID not found")
	}

	// ===== 3. Fetch Data =====
	achievements, err := s.Repo.GetStudentByAchievement(
		context.Background(),
		studentID,
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch achievements")
	}

	// ===== 4. Response =====
	return c.JSON(fiber.Map{
		"message": "Student achievements fetched successfully",
		"total":   len(achievements),
		"data":    achievements,
	})
}

func (s *AchievementService) UploadAttachments(c *fiber.Ctx) error {
	claims := c.Locals("claims")

	fmt.Println("CLAIMS IN SERVICE:", claims) // debugging
	if claims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userClaims := claims.(*middleware.Claims)

	// Hanya admin yang bisa
	if userClaims.Role != "mahasiswa" {
		return fiber.NewError(fiber.StatusForbidden, "Access ditolak, hanya mahasiswa yang boleh mengakses")
	}

	achievementID := c.Params("id")
	if achievementID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Achievement ID is required")
	}

	// ===== 3. Check ownership =====
	achievement, err := s.Repo.FindById(context.Background(), achievementID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Achievment not found")
	}

	if achievement.StudentID != userClaims.UserID {
		return fiber.NewError(fiber.StatusForbidden, "You can only upload to your own achievement")
	}

	// ===== 4. Get file =====
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "File is required")
	}

	// ===== 5. Save file =====
	uploadDir := "./uploads/achievements"
	os.MkdirAll(uploadDir, os.ModePerm)

	ext := filepath.Ext(file.Filename)
	fileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, fileName)

	if err := c.SaveFile(file, filePath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save file")
	}

	// ===== 6. Save metadata =====
	attachment := model.Attachment{
		FileName:   file.Filename,
		FileUrl:    "/uploads/achievements/" + fileName,
		FileType:   file.Header.Get("Content-Type"),
		UploadedAt: time.Now(),
	}

	err = s.Repo.AddAttachment(
		context.Background(),
		achievementID,
		attachment,
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save attachment")
	}

	// ===== 7. Response =====
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Attachment uploaded successfully",
		"data":    attachment,
	})
}
