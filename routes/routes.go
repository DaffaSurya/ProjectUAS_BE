package routes

import (
	"PROJECTUAS_BE/app/service"
	"PROJECTUAS_BE/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authservice *service.AuthService) {
	api := app.Group("/api")

	api.Post("/register", authservice.Register)
	api.Post("/login", authservice.Login)
	api.Get("/Getprofile", middleware.AuthRequired(), authservice.GetProfile)
	api.Post("/logout", middleware.AuthRequired(), authservice.Logout)
}
