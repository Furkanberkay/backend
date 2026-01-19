package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"github.com/skip2/go-qrcode"
)

type TicketService struct {
	eventRepo  models.EventRepository
	repository models.TicketRepository
	logger     *slog.Logger
}

func NewTicketService(repository models.TicketRepository, eventRepo models.EventRepository, logger *slog.Logger) models.TicketService {
	if logger == nil {
		logger = slog.Default()
	}
	return &TicketService{
		repository: repository,
		logger:     logger,
		eventRepo:  eventRepo,
	}
}

func (s *TicketService) GetMany(ctx context.Context, userId uint) ([]*models.Ticket, error) {
	tickets, err := s.repository.GetMany(ctx, userId)
	if err != nil {
		s.logger.ErrorContext(ctx, "tickets_get_many_failed", "error", err)
		return nil, models.InternalError
	}
	return tickets, nil
}

func (s *TicketService) GetOne(ctx context.Context, userId uint, ticketID uint) (*models.Ticket, string, error) {
	// 1. Veriyi Çek
	ticket, err := s.repository.GetOne(ctx, userId, ticketID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			s.logger.DebugContext(ctx, "ticket_not_found", "ticket_id", ticketID)
			return nil, "", models.ErrRecordNotFound
		}

		s.logger.ErrorContext(ctx, "ticket_get_one_failed", "ticket_id", ticketID, "error", err)
		return nil, "", models.InternalError
	}

	qrCodeBytes, err := qrcode.Encode(
		fmt.Sprintf("ticketId:%v,userId:%v", ticketID, userId),
		qrcode.Medium,
		256,
	)
	if err != nil {
		s.logger.Error("qr_generation_failed", "error", err)
		return nil, "", models.InternalError
	}

	qrString := base64.StdEncoding.EncodeToString(qrCodeBytes)

	return ticket, qrString, nil
}

func (s *TicketService) CreateOne(ctx context.Context, userId uint, ticket *models.Ticket) (*models.Ticket, error) {
	_, err := s.eventRepo.GetOne(ctx, ticket.EventID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, models.ErrEventNotFound
		}
		return nil, models.InternalError
	}

	created, err := s.repository.CreateOne(ctx, userId, ticket)
	if err != nil {
		s.logger.ErrorContext(ctx, "ticket_create_failed", "event_id", ticket.EventID, "error", err)
		return nil, models.InternalError
	}
	return created, nil
}

func (s *TicketService) UpdateOne(ctx context.Context, userId uint, ticketID uint, input *models.UpdateTicketInput) (*models.Ticket, error) {

	if input == nil {
		return nil, models.NewValidationError("update data cannot be empty")
	}

	if input.Price != nil && *input.Price < 0 {
		return nil, models.NewValidationError("price cannot be negative")
	}

	if input.Status != nil {
		validStatuses := map[string]bool{"sold": true, "available": true, "reserved": true}
		if !validStatuses[*input.Status] {
			return nil, models.NewValidationError("invalid status: " + *input.Status)
		}
	}

	updated, err := s.repository.UpdateOne(ctx, userId, ticketID, input)

	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, models.ErrRecordNotFound
		}
		s.logger.ErrorContext(ctx, "ticket_update_failed", "ticket_id", ticketID, "error", err)
		return nil, models.InternalError
	}

	return updated, nil
}

func (s *TicketService) ValidateEntry(ctx context.Context, userId uint, ticketID uint) (*models.Ticket, error) {
	ticket, err := s.repository.GetOne(ctx, userId, ticketID)
	if err != nil {
		return nil, err
	}

	if ticket.Entered {
		return nil, models.ErrTicketAlreadyUsed
	}

	entered := true
	input := &models.UpdateTicketInput{
		Entered: &entered,
	}

	return s.repository.UpdateOne(ctx, userId, ticketID, input)
}
