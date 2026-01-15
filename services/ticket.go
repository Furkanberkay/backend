package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
)

type TicketService struct {
	repository models.TicketRepository
	logger     *slog.Logger
}

func NewTicketService(repository models.TicketRepository, logger *slog.Logger) models.TicketService {
	if logger == nil {
		logger = slog.Default()
	}
	return &TicketService{
		repository: repository,
		logger:     logger,
	}
}

func (s *TicketService) GetMany(ctx context.Context) ([]*models.Ticket, error) {
	tickets, err := s.repository.GetMany(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "tickets_get_many_failed", "error", err)
		return nil, models.InternalError
	}
	return tickets, nil
}

func (s *TicketService) GetOne(ctx context.Context, ticketID uint) (*models.Ticket, error) {
	ticket, err := s.repository.GetOne(ctx, ticketID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			s.logger.DebugContext(ctx, "ticket_not_found", "ticket_id", ticketID)
			return nil, models.ErrRecordNotFound
		}

		s.logger.ErrorContext(ctx, "ticket_get_one_failed", "ticket_id", ticketID, "error", err)
		return nil, models.InternalError
	}
	return ticket, nil
}

func (s *TicketService) CreateOne(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {

	created, err := s.repository.CreateOne(ctx, ticket)
	if err != nil {
		s.logger.ErrorContext(ctx, "ticket_create_failed", "event_id", ticket.EventID, "error", err)
		return nil, models.InternalError
	}
	return created, nil
}

func (s *TicketService) UpdateOne(ctx context.Context, ticketID uint, input *models.UpdateTicketInput) (*models.Ticket, error) {

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

	updated, err := s.repository.UpdateOne(ctx, ticketID, input)

	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, models.ErrRecordNotFound
		}
		s.logger.ErrorContext(ctx, "ticket_update_failed", "ticket_id", ticketID, "error", err)
		return nil, models.InternalError
	}

	return updated, nil
}

func (s *TicketService) ValidateEntry(ctx context.Context, ticketID uint) (*models.Ticket, error) {
	ticket, err := s.repository.GetOne(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	if ticket.Entered {
		return nil, models.NewValidationError("ticket already used")
	}

	entered := true

	input := &models.UpdateTicketInput{
		Entered: &entered,
	}

	return s.repository.UpdateOne(ctx, ticketID, input)
}
