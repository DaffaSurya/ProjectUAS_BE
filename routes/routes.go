package routes

import (
	"PROJECTUAS_BE/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authservice *service.AuthService) {
	// Grouping routes
	// admin := app.Group("/admin", middleware.AuthRequired(), middleware.RequireRole("admin"))
	// mahasiswa := app.Group("/mahasiswa", middleware.AuthRequired(), middleware.RequireRole("mahasiswa"))
	// dosen := app.Group("/dosen", middleware.AuthRequired(), middleware.RequireRole("dosen"))
	app.Post("/login", authservice.Login)
	// ============================
	// ADMIN ROUTES
	// ============================
	// admin.Get("/", func(c *fiber.Ctx) error {
	// 	return c.JSON(fiber.Map{
	// 		"message": "Admin dashboard",
	// 	})
	// })

	// ============================
	// MAHASISWA ROUTES
	// ============================
	// mahasiswa.Post("/prestasi", func(c *fiber.Ctx) error {
	// 	return c.JSON(fiber.Map{
	// 		"message": "Prestasi berhasil dilaporkan",
	// 	})
	// })

	// ============================
	// DOSEN ROUTES
	// ============================
	// 	dosen.Post("/verify", func(c *fiber.Ctx) error {
	// 		return c.JSON(fiber.Map{
	// 			"message": "Prestasi berhasil diverifikasi",
	// 		})
	// 	})
}
