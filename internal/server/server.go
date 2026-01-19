package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rafli2460/culinary-blog-api/internal/domain"
	"github.com/rs/zerolog"
)

type App struct {
	Fiber    *fiber.App
	Logger   *zerolog.Logger
	Ds       *Datasources
	Services *Services
	config   map[string]string
}

type Services struct {
	User domain.User
}

type Datasources struct {
	WriterDB *sqlx.DB `json:"writer-db"`
	ReaderDB *sqlx.DB `json:"reader-db"`
}
