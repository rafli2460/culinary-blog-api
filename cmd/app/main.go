package main

import (
	"context"

	"github.com/rafli2460/culinary-blog-api/internal/config"
	"github.com/rafli2460/culinary-blog-api/internal/modules"
	"github.com/rafli2460/culinary-blog-api/internal/platform/database"
	"github.com/rafli2460/culinary-blog-api/internal/server"
)

func main() {
	// TODO: Load config from env or file
	conf := map[string]string{
		config.DbHostWriter: "localhost",
		config.DbHostReader: "localhost",
		config.DbPort:       "3306",
		config.DbUser:       "root",
		config.DbPass:       "password",
		config.DbName:       "culinary_blog",
		config.DbDialeg:     "mysql",
		config.ServerPort:   "8080",
	}

	app := server.New(conf)

	// Initialize Database
	app.Ds = database.Init(conf)

	// Initialize Modules (Services & Handlers)
	app.Services = modules.Init(context.Background(), app)

	// Map Routes
	app.MapRoutes()

	// Start Server
	app.Logger.Info().Msg("Starting server on port " + conf[config.ServerPort])
	if err := app.Fiber.Listen(":" + conf[config.ServerPort]); err != nil {
		app.Logger.Fatal().Err(err).Msg("Server failed to start")
	}
}