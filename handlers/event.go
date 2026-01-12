package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	repository models.EventService
}

func NewEventHandler(router fiber.Router, repository models.EventRepository) *EventHandler {
	handler := &EventHandler{
		repository: repository,
	}

	router.Get("/", handler.GetMany)
	router.Post("/", handler.CreateOne)
	router.Get("/:eventId", handler.GetOne)
	return handler
}

func (h *EventHandler) GetMany(ctx *fiber.Ctx) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	events, err := h.repository.GetMany(timeoutCtx)
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    events,
	})
}

func (h *EventHandler) GetOne(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}

func (h *EventHandler) CreateOne(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"status": "success",
	})
}
