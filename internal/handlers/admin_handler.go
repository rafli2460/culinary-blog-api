package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/rafli2460/culinary-blog-api/internal/models"
	"github.com/rafli2460/culinary-blog-api/internal/service"
	"github.com/rs/zerolog/log"
)

type AdminHandler struct {
	userService service.UserService
}

func NewAdminHandler(userService service.UserService) *AdminHandler {
	return &AdminHandler{userService: userService}
}

func getAdminID(c fiber.Ctx) int {
	userIDVal := c.Locals("user_id")
	if idFloat, ok := userIDVal.(float64); ok {
		return int(idFloat)
	}

	if idInt, ok := userIDVal.(int); ok {
		return idInt
	}
	return 0
}

func (h *AdminHandler) GetUsers(c fiber.Ctx) error {
	search := c.Query("search")

	users, err := h.userService.GetAllUsers(c.Context(), search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to retrieve data",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}

func (h *AdminHandler) GetStats(c fiber.Ctx) error {
	stats, err := h.userService.GetStats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to retrieve user data",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   stats,
	})
}

func (h *AdminHandler) UpdateRole(c fiber.Ctx) error {
	targetID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid user ID",
		})
	}

	var req models.UpdateRoleRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid data format",
		})
	}

	currentAdminID := getAdminID(c)

	err = h.userService.UpdateRole(c.Context(), targetID, currentAdminID, req.Role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": err.Error(),
		})
	}

	log.Info().Int("admin_id", currentAdminID).Int("target_id", targetID).Str("new_role", req.Role).Msg("Role successfully updated")

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User role successfully updated",
	})
}

func (h *AdminHandler) DeleteUser(c fiber.Ctx) error {
	targetID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid user ID",
		})
	}

	currentAdminID := getAdminID(c)

	err = h.userService.DeleteUser(c.Context(), targetID, currentAdminID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": err.Error(),
		})
	}

	log.Info().Int("admin_id", currentAdminID).Int("target_id", targetID).Msg("User successfully deleted")

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User successfully deleted",
	})
}
