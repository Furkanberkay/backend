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

type TicketHandler struct {
	service  models.TicketService
	validate *validator.Validate
}

func NewTicketHandler(validate *validator.Validate, router fiber.Router, service models.TicketService) {
	handler := &TicketHandler{service: service, validate: validate}

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
		return httpx.MapErrorToResponse(c, err)
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
		return httpx.MapErrorToResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   ticket,
	})
}

func (h *TicketHandler) CreateOne(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	var input dto.CreateTicketInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid body"})
	}

	if err := h.validate.Struct(input); err != nil {
		return httpx.MapErrorToResponse(c, err)
	}
	ticket := new(models.Ticket)
	ticket.EventID = input.EventID

	createdTicket, err := h.service.CreateOne(ctx, ticket)
	if err != nil {
		return httpx.MapErrorToResponse(c, err)
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
		return httpx.MapErrorToResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   updatedTicket,
	})
}

func (h *TicketHandler) ValidateEntry(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	req := new(dto.ValidateTicket)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid body",
		})
	}

	ticket, err := h.service.ValidateEntry(ctx, req.TicketID)
	if err != nil {
		return httpx.MapErrorToResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Entry approved. Welcome!",
		"data":    ticket,
	})
}
