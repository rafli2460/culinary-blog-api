package middleware

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rafli2460/culinary-blog-api/pkg/response"
	"github.com/rs/zerolog/log"
)

func Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenString := c.Cookies("jwt_token")

		if tokenString == "" {
			log.Warn().Msg("Access Denied: token not found")
			return response.Error(c, fiber.StatusUnauthorized, "access denied")
		}

		secretKey := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			log.Warn().Err(err).Msg("token invalid: token expired or not valid")
			return response.Error(c, fiber.StatusUnauthorized, "session is not valid or has ended. Please re-login")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("user_id", claims["user_id"])
		}

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenString := c.Cookies("jwt_token")

		if tokenString == "" {
			log.Warn().Msg("Access Denied: token not found")
			return response.Error(c, fiber.StatusUnauthorized, "access denied")
		}

		secretKey := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			log.Warn().Err(err).Msg("token invalid: token expired or not valid")
			return response.Error(c, fiber.StatusUnauthorized, "session is not valid or has ended. Please re-login")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			role, roleExists := claims["role"].(string)
			if !roleExists || role != "admin" {
				log.Warn().Interface("user_id", claims["user_id"]).Msg("admin access attempts by non-admins")
				return response.Error(c, fiber.StatusForbidden, "Access prohibited. You do not have admin permission.")
			}
			c.Locals("user_id", claims["user_id"])
			c.Locals("role", role)
		}

		return c.Next()
	}
}
