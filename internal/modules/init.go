package modules

import (
	"context"

	"github.com/rafli2460/culinary-blog-api/internal/modules/users"
	"github.com/rafli2460/culinary-blog-api/internal/server"
)

func Init(ctx context.Context, app *server.App) *server.Services {
	userService, userHandler := users.Init(ctx, app)

	srv := &server.Services{
		User:        userService,
		UserHandler: userHandler,
	}

	return srv
}