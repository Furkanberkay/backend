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
	handler := &AuthHandler{
		service:  service,
		validate: validate,
	}

	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)
	router.Post("/refresh", handler.Refresh)

	return handler
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	input := new(dto.RegisterUserRequest)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid request body"})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	userModel := &models.User{
		Name:     input.Name,
		Surname:  input.Surname,
		Email:    input.Email,
		Birthday: input.Birthday,
	}

	createdUser, err := h.service.Register(ctx, userModel, input.Password)
	if err != nil {
		if err == models.ErrUserAlreadyExist {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "User already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   dto.ToUserResponse(createdUser),
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	input := new(dto.LoginRequest)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid body"})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	creds := &models.AuthCredentials{
		Email:    input.Email,
		Password: input.Password,
	}

	authData, err := h.service.Login(ctx, creds)
	if err != nil {
		if err == models.ErrInvalidCredentials || err == models.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid email or password"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	response := dto.LoginResponse{
		AccessToken:  authData.AccessToken,
		RefreshToken: authData.RefreshToken,
		User:         dto.ToUserResponse(authData.User),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	input := new(dto.RefreshTokenRequest)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid body"})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	authData, err := h.service.RefreshToken(ctx, input.RefreshToken)
	if err != nil {
		if err == models.ErrInvalidToken || err == models.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid or expired refresh token"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	response := dto.LoginResponse{
		AccessToken:  authData.AccessToken,
		RefreshToken: authData.RefreshToken,
		User:         dto.ToUserResponse(authData.User),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}
