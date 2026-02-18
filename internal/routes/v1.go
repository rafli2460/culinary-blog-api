package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/rafli2460/culinary-blog-api/internal/handlers"
	"github.com/rafli2460/culinary-blog-api/internal/middleware"
)

func InitRoutes(app *fiber.App,
	authHandler *handlers.AuthHandler,
	adminHandler *handlers.AdminHandler,
	postHandler *handlers.PostHandler) {

	app.Get("/uploads/*", static.New("./uploads"))
	api := app.Group("/v1")

	api.Get("/posts/:id", postHandler.GetPost)
	api.Get("/posts", postHandler.GetAllPosts)

	api.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "system is healthy",
		})
	})

	// AUTH
	auth := api.Group("/auth")

	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	// ADMIN
	admin := api.Group("/admin", middleware.AdminOnly())

	admin.Get("/users/stats", adminHandler.GetStats)
	admin.Get("/users", adminHandler.GetUsers)

	admin.Put("/users/:id/role", adminHandler.UpdateRole)
	admin.Delete("/users/:id", adminHandler.DeleteUser)

	// POST
	posts := api.Group("/post", middleware.Protected())
	posts.Post("/", postHandler.CreatePost)
	posts.Delete("/:id", postHandler.DeletePost)
	posts.Put("/:id", postHandler.UpdatePost)

}
