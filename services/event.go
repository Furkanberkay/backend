package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
)

type EventService struct {
	repository models.EventRepository
	logger     *slog.Logger
}

func NewEventService(repository models.EventRepository, logger *slog.Logger) models.EventService {
	return &EventService{repository: repository, logger: logger}
}

func (s *EventService) GetMany(ctx context.Context) ([]*models.Event, error) {
	events, err := s.repository.GetMany(ctx)
	if err != nil {
		return nil, models.InternalError
	}
	return events, nil
}

func (s *EventService) GetOne(ctx context.Context, eventId uint) (*models.Event, error) {
	event, err := s.repository.GetOne(ctx, eventId)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, models.ErrRecordNotFound
		}
		return nil, models.InternalError
	}
	return event, nil
}

func (s *EventService) CreateOne(ctx context.Context, event *models.Event) (*models.Event, error) {
	created, err := s.repository.CreateOne(ctx, event)
	if err != nil {
		return nil, models.InternalError
	}
	return created, nil
}

func (s *EventService) UpdateOne(ctx context.Context, eventId uint, event *models.Event) (*models.Event, error) {
	updated, err := s.repository.UpdateOne(ctx, eventId, event)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, models.ErrRecordNotFound
		}
		return nil, models.InternalError
	}
	return updated, nil
}

func (s *EventService) PatchOne(ctx context.Context, eventId uint, input *models.EventUpdateInput) (*models.Event, error) {
	current, err := s.repository.GetOne(ctx, eventId)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, models.ErrRecordNotFound
		}
		s.logger.ErrorContext(ctx, "patch_event_get_failed", "event_id", eventId, "error", err)
		return nil, models.InternalError
	}

	if input == nil {
		return current, nil
	}

	isChanged := false

	if input.Name != nil {
		if len(*input.Name) < 3 {
			return nil, models.NewValidationError("event name too short")
		}
		if current.Name != *input.Name {
			current.Name = *input.Name
			isChanged = true
		}
	}

	if input.Location != nil {
		if current.Location != *input.Location {
			current.Location = *input.Location
			isChanged = true
		}
	}

	if input.Date != nil {
		if input.Date.Before(time.Now()) {
			return nil, models.NewValidationError("event date cannot be in the past")
		}
		if !current.Date.Equal(*input.Date) {
			current.Date = *input.Date
			isChanged = true
		}
	}

	if !isChanged {
		return current, nil
	}

	updated, err := s.repository.UpdateOne(ctx, eventId, current)
	if err != nil {
		s.logger.ErrorContext(ctx, "patch_event_update_failed", "event_id", eventId, "error", err)
		return nil, models.InternalError
	}

	s.logger.InfoContext(ctx, "event_patched_successfully", "event_id", eventId)
	return updated, nil
}
func (s *EventService) DeleteOne(ctx context.Context, eventId uint) error {
	err := s.repository.DeleteOne(ctx, eventId)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return models.ErrRecordNotFound
		}
		return models.InternalError
	}
	return nil
}
