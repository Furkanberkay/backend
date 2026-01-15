package handlers

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	service models.TicketService
}

func NewTicketHandler(router fiber.Router, service models.TicketService) {
	handler := &TicketHandler{service: service}

	router.Get("/", handler.GetMany)
	router.Post("/", handler.CreateOne)
	router.Get("/:id", handler.GetOne)

	router.Patch("/:id", handler.UpdateOne)

	router.Post("/validate", handler.ValidateEntry)
}

func (h *TicketHandler) GetMany(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	tickets, err := h.service.GetMany(ctx)
	if err != nil {
		return mapErrorToResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   tickets,
	})
}

func (h *TicketHandler) GetOne(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid ticket ID",
		})
	}

	ticket, err := h.service.GetOne(ctx, uint(id))
	if err != nil {
		return mapErrorToResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   ticket,
	})
}

func (h *TicketHandler) CreateOne(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	ticket := new(models.Ticket)
	if err := c.BodyParser(ticket); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	createdTicket, err := h.service.CreateOne(ctx, ticket)
	if err != nil {
		return mapErrorToResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   createdTicket,
	})
}

func (h *TicketHandler) UpdateOne(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid ticket ID",
		})
	}

	input := new(models.UpdateTicketInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	updatedTicket, err := h.service.UpdateOne(ctx, uint(id), input)
	if err != nil {
		return mapErrorToResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   updatedTicket,
	})
}

func (h *TicketHandler) ValidateEntry(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	type EntryRequest struct {
		TicketID uint `json:"ticket_id"`
	}
	req := new(EntryRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid body",
		})
	}

	ticket, err := h.service.ValidateEntry(ctx, req.TicketID)
	if err != nil {
		return mapErrorToResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Entry approved. Welcome!",
		"data":    ticket,
	})
}

func mapErrorToResponse(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, models.ErrRecordNotFound) || errors.Is(err, models.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Resource not found",
		})
	case errors.Is(err, models.ErrValidation):
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
