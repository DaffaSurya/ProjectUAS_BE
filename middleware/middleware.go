package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing token",
			})
		}

		parts := strings.Split(authHeader, " ")
		// Contoh: "Bearer tokenxxxxx"
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token format",
			})
		}

		tokenStr := parts[1]

		// CEK TOKEN BLACKLIST
		if IsTokenBlacklisted(tokenStr) { // tidak perlu prefix middleware.
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token has been revoked",
			})
		}

		// Parse JWT
		claims, err := ParseToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		// DEBUG: memastikan claims terbaca
		fmt.Println("=== MIDDLEWARE CLAIMS ===")
		fmt.Printf("Email: %s\n", claims.Email) // debugging part
		fmt.Println("==========================")

		// Simpan ke fiber locals
		c.Locals("claims", claims)      // bentuk struct claims
		c.Locals("email", claims.Email) // bisa dipakai jika butuh email saja

		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(*Claims)

		for _, r := range roles {
			if claims.Role == r {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden, role not allowed",
		})
	}
}

func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(*Claims)

		for _, p := range claims.Permissions {
			if p == permission {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "permission denied",
		})
	}
}
