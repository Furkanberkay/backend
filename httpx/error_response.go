package httpx

import (
	"errors"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"github.com/gofiber/fiber/v2"
)

func MapErrorToResponse(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, models.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Resource not found",
		})

	case errors.Is(err, models.ErrValidation):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})

	case errors.Is(err, models.ErrTicketAlreadyUsed):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})

	case errors.Is(err, models.ErrEventNotFound):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	case errors.Is(err, models.ErrUserAlreadyExist):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
		})
	}
}
