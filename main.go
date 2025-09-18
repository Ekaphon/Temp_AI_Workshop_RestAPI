// @title Temp AI Workshop REST API
// @version 1.0
// @description Simple API for authentication using Fiber, GORM and SQLite
// @host localhost:3000
// @BasePath /
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	swag "github.com/gofiber/swagger"
	_ "temp-ai-restapi/docs"
)

func main() {
	InitDB()
	app := fiber.New()

	// Auth routes
	api := app.Group("/auth")
	api.Post("/register", Register)
	api.Post("/login", Login)

	// Protected example
	app.Get("/profile", RequireAuth, func(c *fiber.Ctx) error {
		u := c.Locals("user").(*User)
		return c.JSON(fiber.Map{"email": u.Email, "id": u.ID, "first_name": u.FirstName, "last_name": u.LastName, "phone": u.Phone, "member_level": u.MemberLevel, "points": u.Points})
	})
	app.Put("/profile", RequireAuth, UpdateProfile)

	// Health
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "hello world"})
	})

	// Swagger UI
	app.Get("/swagger/*", swag.HandlerDefault)

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
