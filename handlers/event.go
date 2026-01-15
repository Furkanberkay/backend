package handlers

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	service models.EventService
}

func NewEventHandler(router fiber.Router, service models.EventService) *EventHandler {
	handler := &EventHandler{service: service}

	router.Get("/", handler.GetMany)
	router.Post("/", handler.CreateOne)
	router.Get("/:eventId", handler.GetOne)
	router.Put("/:eventId", handler.UpdateOne)
	router.Patch("/:eventId", handler.PatchOne)
	router.Delete("/:eventId", handler.DeleteOne)

	return handler
}

func (h *EventHandler) GetMany(ctx *fiber.Ctx) error {
	timeoutCtx, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	events, err := h.service.GetMany(timeoutCtx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": models.InternalError.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    events,
	})
}

func (h *EventHandler) GetOne(ctx *fiber.Ctx) error {
	eventId, err := strconv.Atoi(ctx.Params("eventId"))
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
		if errors.Is(err, models.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": models.ErrRecordNotFound.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": models.InternalError.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    event,
	})
}

func (h *EventHandler) CreateOne(ctx *fiber.Ctx) error {
	event := new(models.Event)
	if err := ctx.BodyParser(event); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	createdEvent, err := h.service.CreateOne(timeoutCtx, event)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": models.InternalError.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Event created successfully",
		"data":    createdEvent,
	})
}

func (h *EventHandler) UpdateOne(ctx *fiber.Ctx) error {
	eventId, err := strconv.Atoi(ctx.Params("eventId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	event := new(models.Event)
	if err := ctx.BodyParser(event); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	updatedEvent, err := h.service.UpdateOne(timeoutCtx, uint(eventId), event)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": models.ErrRecordNotFound.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": models.InternalError.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Event updated successfully",
		"data":    updatedEvent,
	})
}

func (h *EventHandler) PatchOne(ctx *fiber.Ctx) error {
	eventId, err := strconv.Atoi(ctx.Params("eventId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	patch := new(models.EventPatchDTO)
	if err := ctx.BodyParser(patch); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	patchedEvent, err := h.service.PatchOne(timeoutCtx, uint(eventId), patch)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": models.ErrRecordNotFound.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": models.InternalError.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Event patched successfully",
		"data":    patchedEvent,
	})
}

func (h *EventHandler) DeleteOne(ctx *fiber.Ctx) error {
	eventId, err := strconv.Atoi(ctx.Params("eventId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.service.DeleteOne(timeoutCtx, uint(eventId)); err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": models.ErrRecordNotFound.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": models.InternalError.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Event deleted successfully",
	})
}
