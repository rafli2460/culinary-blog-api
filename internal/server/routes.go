package server

import "github.com/gofiber/fiber/v2"

func (a *App) MapRoutes() {
	v1 := a.Fiber.Group("/api/v1")

	// Users Routes
	users := v1.Group("/users")
	users.Post("/register", a.Services.UserHandler.Register)
	users.Post("/login", a.Services.UserHandler.Login)
	users.Put("/:id/role", a.Services.UserHandler.UpdateRole)
	users.Delete("/:id", a.Services.UserHandler.DeleteUser)

	// Health check
	a.Fiber.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
