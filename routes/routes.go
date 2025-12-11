package routes

import (
	"PROJECTUAS_BE/app/service"
	"PROJECTUAS_BE/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authservice *service.AuthService, Studentservice *service.Studentservice, AchieveService *service.AchievementService) {
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

	api.Use(middleware.AuthRequired()) // melindungi agar hanya admin yang bisa mengakses
	api.Post("/achievements", AchieveService.CreateAchievements)
	api.Put("/achievements/:id", AchieveService.UpdateAchievement)
	api.Delete("/achievements/:id", AchieveService.DeleteAchievement)

	// students
	api.Get("/student/:id", Studentservice.GetStudent)
	api.Get("/student", Studentservice.GetAllStudents)

}
