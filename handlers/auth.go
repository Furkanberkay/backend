package handlers

import (
	"context"
	"time"

	"github.com/Furkanberkay/ticket-booking-project-v1/dto"
	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service  models.AuthService
	validate *validator.Validate
}

func NewAuthHandler(validate *validator.Validate, router fiber.Router, service models.AuthService) *AuthHandler {
	handler := &AuthHandler{service: service, validate: validate}

	router.Get("/", handler.GetMany)

	return handler
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	input := new(dto.RegisterUserRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid body"})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	user := models.User{
		Name:     input.Name,
		Surname:  input.Surname,
		Email:    input.Email,
		Birthday: input.Birthday,
	}

	h.service.Register(ctx, &user, input.Password)

}
