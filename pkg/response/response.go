package response

import (
	"github.com/gofiber/fiber/v3"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(c fiber.Ctx, statusCode int, message string, data interface{}, meta interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Status:  "error",
		Message: message,
	})
}
