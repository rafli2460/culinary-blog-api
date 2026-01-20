package server

import (
	"os"

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

type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	UpdateRole(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type Services struct {
	User        domain.UserService
	UserHandler UserHandler
}

type Datasources struct {
	WriterDB *sqlx.DB `json:"writer-db"`
	ReaderDB *sqlx.DB `json:"reader-db"`
}

func New(conf map[string]string) *App {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	return &App{
		Fiber:  fiber.New(),
		Logger: &logger,
		config: conf,
	}
}