package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rafli2460/culinary-blog-api/internal/handlers"
	"github.com/rafli2460/culinary-blog-api/internal/middleware"
)

func InitRoutes(app *fiber.App, authHandler *handlers.AuthHandler, adminHandler *handlers.AdminHandler) {
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

	admin := api.Group("/admin", middleware.AdminOnly())

	admin.Get("/users/stats", adminHandler.GetStats)
	admin.Get("/users", adminHandler.GetUsers)

	admin.Put("/users/:id/role", adminHandler.UpdateRole)
	admin.Delete("/users/:id", adminHandler.DeleteUser)

}
