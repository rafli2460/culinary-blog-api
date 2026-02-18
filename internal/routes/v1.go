package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rafli2460/culinary-blog-api/internal/handlers"
)

func InitRoutes(app *fiber.App, authHandler *handlers.AuthHandler) {
	api := app.Group("/api")

	api.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "system is healthy",
		})
	})

	auth := api.Group("/auth")

	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

}
