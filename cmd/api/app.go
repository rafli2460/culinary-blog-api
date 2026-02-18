package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/rafli2460/culinary-blog-api/internal/config"
	"github.com/rafli2460/culinary-blog-api/internal/handlers"
	"github.com/rafli2460/culinary-blog-api/internal/repository"
	"github.com/rafli2460/culinary-blog-api/internal/routes"
	"github.com/rafli2460/culinary-blog-api/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load()
	if err != nil {
		log.Warn().Msg(".env file not found, using system environment")
	}

	db := config.InitDB()
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)

	authHandler := handlers.NewAuthHandler(userService)
	adminHandler := handlers.NewAdminHandler(userService)
	postHandler := handlers.NewPostService(postService)
	app := fiber.New()

	routes.InitRoutes(app, authHandler, adminHandler, postHandler)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3000"
	}

	log.Info().Msgf("Running server on port %s", appPort)
	err = app.Listen(":"+appPort, fiber.ListenConfig{
		DisableStartupMessage: true,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("server failed to run")
	}
}
