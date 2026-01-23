package responses

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rafli2460/culinary-blog-api/internal/config"
)

func HttpError(c *fiber.Ctx, err error) error {
	var errResponse *ErrorResponse
	if errors.As(err, &errResponse) {
		c.Status(errResponse.Status)
	}

	if errResponse == nil {
		errResponse = &ErrorResponse{}
	}

	if os.Getenv(config.ServerEnv) == config.EnvDevelopment {
		errResponse.Debug = errResponse.Error()
	}

	c.Append("Access-Control-Allow-Origin", "*")

	//app.Logger.Error().Stack().
	//	Str("Method", c.Method()).
	//	Str("Path", c.Path()).
	//	Int("Status", errResponse.Status).
	//	Err(err).Msg(errResponse.Message)

	_ = c.JSON(errResponse)

	return nil
}

func HttpSuccess(c *fiber.Ctx, message string, data interface{}) (err error) {
	response := Response{}
	response.Status = fiber.StatusOK
	response.Message = message
	response.Data = data

	c.Append("Access-Control-Allow-Origin", "*")

	//app.Logger.Log().
	//	Str("Method", c.Method()).
	//	Str("Path", c.Path()).
	//	Str("Status", response.Status).
	//	Msg(response.Message)

	_ = c.JSON(response)

	return nil
}
