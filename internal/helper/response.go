package helper

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	message string
}

func ResponseWithError(c *fiber.Ctx, err error) error {
	errorResponse := ErrorResponse{
		message: err.Error(),
	}
	return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
}

func ResponseWithErrorMessage(c *fiber.Ctx, msg string) error {
	errorResponse := ErrorResponse{
		message: msg,
	}
	return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
}
