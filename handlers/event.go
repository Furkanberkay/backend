package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/Furkanberkay/ticket-booking-project-v1/dto"
	"github.com/Furkanberkay/ticket-booking-project-v1/httpx"
	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	service  models.EventService
	validate *validator.Validate
}

func NewEventHandler(validate *validator.Validate, router fiber.Router, service models.EventService) *EventHandler {
	handler := &EventHandler{service: service, validate: validate}

	router.Get("/", handler.GetMany)
	router.Post("/", handler.CreateOne)
	router.Get("/:id", handler.GetOne)
	router.Put("/:id", handler.UpdateOne)
	router.Patch("/:id", handler.PatchOne)
	router.Delete("/:id", handler.DeleteOne)

	return handler
}

func (h *EventHandler) GetMany(ctx *fiber.Ctx) error {
	timeoutCtx, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	events, err := h.service.GetMany(timeoutCtx)
	if err != nil {
		return httpx.MapErrorToResponse(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    events,
	})
}

func (h *EventHandler) GetOne(ctx *fiber.Ctx) error {
	eventId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil || eventId <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "invalid eventId",
		})
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	event, err := h.service.GetOne(timeoutCtx, uint(eventId))
	if err != nil {
		return httpx.MapErrorToResponse(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    event,
	})
}

func (h *EventHandler) CreateOne(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	input := new(dto.CreateEventInput)
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

	event := &models.Event{
		Name:     input.Name,
		Location: input.Location,
		Date:     input.Date,
	}

	createdEvent, err := h.service.CreateOne(ctx, event)
	if err != nil {
		return httpx.MapErrorToResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   createdEvent,
	})
}

func (h *EventHandler) UpdateOne(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid ID"})
	}

	input := new(dto.UpdateEventInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid body"})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	event := &models.Event{
		Name:     input.Name,
		Location: input.Location,
		Date:     input.Date,
	}

	updatedEvent, err := h.service.UpdateOne(ctx, uint(id), event)
	if err != nil {
		return httpx.MapErrorToResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   updatedEvent,
	})
}

func (h *EventHandler) PatchOne(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid ID"})
	}

	input := new(dto.EventPatchRequest)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid body"})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	event := models.EventUpdateInput{
		Name:     input.Name,
		Location: input.Location,
		Date:     input.Date,
	}

	patchedEvent, err := h.service.PatchOne(ctx, uint(id), &event)
	if err != nil {
		return httpx.MapErrorToResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   patchedEvent,
	})
}

func (h *EventHandler) DeleteOne(ctx *fiber.Ctx) error {
	eventId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.service.DeleteOne(timeoutCtx, uint(eventId)); err != nil {
		return httpx.MapErrorToResponse(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Event deleted successfully",
	})
}
