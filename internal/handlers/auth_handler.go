package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rafli2460/culinary-blog-api/internal/models"
	"github.com/rafli2460/culinary-blog-api/internal/service"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req models.RegisterRequest

	if err := c.Bind().Body(&req); err != nil {
		log.Warn().Err(err).Msg("failed to parse request registration body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Format is invalid",
		})
	}

	log.Info().Str("Username", req.Username).Msg("User registration success")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "Registration successful, please login",
	})
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "invalid format",
		})
	}

	token, err := h.userService.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Expires:  time.Now().Add(12 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	log.Info().Str("username", req.Username).Msg("User Login Success")

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
	})
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	log.Info().Msg("User logout success")

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logout successful",
	})
}
