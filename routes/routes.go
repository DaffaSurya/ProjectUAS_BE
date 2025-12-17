package routes

import (
	"PROJECTUAS_BE/app/service"
	"PROJECTUAS_BE/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, Userservice *service.UserService, Studentservice *service.Studentservice, AchieveService *service.AchievementService, LectureService *service.LecturesService, ReportService *service.ReportService, AuthService *service.AuthService) {
	api := app.Group("/api")

	// authentication Route
	api.Post("/login", AuthService.Login)
	api.Get("/Getprofile", middleware.AuthRequired(), AuthService.GetProfile)
	api.Post("/logout", middleware.AuthRequired(), AuthService.Logout)
	// authentication route

	// users route
	api.Use(middleware.AuthRequired()) // melindungi agar hanya admin yang bisa mengakses
	api.Get("/users", Userservice.GetAllUsers)
	api.Get("/users/:id", Userservice.GetUsersByID)
	api.Post("/users", Userservice.CreateUser)
	api.Put("/users/:id", Userservice.UpdateUserByID)
	api.Delete("/users/:id", Userservice.DeleteUserByID)

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
