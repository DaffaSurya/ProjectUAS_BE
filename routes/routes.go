package routes

import (
	"PROJECTUAS_BE/app/service"
	"PROJECTUAS_BE/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authservice *service.AuthService, Studentservice *service.Studentservice, AchieveService *service.AchievementService, LectureService *service.LecturesService, ReportService *service.ReportService) {
	api := app.Group("/api")

	// authentication Route
	api.Post("/login", authservice.Login)
	api.Get("/Getprofile", middleware.AuthRequired(), authservice.GetProfile)
	api.Post("/logout", middleware.AuthRequired(), authservice.Logout)
	// authentication route

	// users route
	api.Use(middleware.AuthRequired()) // melindungi agar hanya admin yang bisa mengakses
	api.Get("/users", authservice.GetAllUsers)
	api.Get("/users/:id", authservice.GetUsersByID)
	api.Post("/users", authservice.CreateUser)
	api.Put("/users/:id", authservice.UpdateUserByID)
	api.Delete("/users/:id", authservice.DeleteUserByID)

	// achievement
	api.Get("/achievements", AchieveService.GetAllAchievements)
	api.Get("/achievements/:id", AchieveService.GetAchievementsByID)

	api.Use(middleware.AuthRequired()) // melindungi agar hanya mahasiswa yang bisa mengakses
	api.Post("/achievements", AchieveService.CreateAchievements)
	api.Put("/achievements/:id", AchieveService.UpdateAchievement)
	api.Delete("/achievements/:id", AchieveService.DeleteAchievement)
	api.Post("/achievements/:id/submit", Studentservice.SubmitAchievement)

	api.Use(middleware.AuthRequired()) // melindungi agar hanya dosen wali yang bisa mengakses
	api.Post("/achievements/:id/verify", LectureService.VerifyAchievement)
	api.Post("/achievements/:id/Reject", LectureService.RejectAchievement)

	api.Post("/achievements/:id/history", LectureService.GetHistory)

	api.Post("/achievements/:id/Attachment", AchieveService.UploadAttachments)

	// students and Lecturers
	api.Get("/student/:id", Studentservice.GetStudent)
	api.Get("/student", Studentservice.GetAllStudents)

	api.Use(middleware.AuthRequired())
	api.Get("/student/:id/achievement", AchieveService.GetStudentAchievements)
	api.Put("student/:id/advisor", Studentservice.UpdateAdvisor)
	// lectures
	api.Get("/lecturers", LectureService.GetLectures)
	api.Get("/lecturers/:id/advisees", LectureService.Getadvisees)

	// report and analytics
	api.Use(middleware.AuthRequired())
	api.Get("/reports/statics", ReportService.GetStatics)
	api.Get("/reports/student/:id", ReportService.GetStudentReport)
}
