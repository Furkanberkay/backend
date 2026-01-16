package services

import (
	"context"
	"errors"

	"github.com/Furkanberkay/ticket-booking-project-v1/dto"
	"github.com/Furkanberkay/ticket-booking-project-v1/models"
)

type EventService struct {
	repository models.EventRepository
}

func NewEventService(repository models.EventRepository) models.EventService {
	return &EventService{repository: repository}
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

func (s *EventService) PatchOne(ctx context.Context, eventId uint, patch *dto.EventPatchRequest) (*models.Event, error) {
	current, err := s.repository.GetOne(ctx, eventId)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, models.ErrRecordNotFound
		}
		return nil, models.InternalError
	}

	if patch != nil {
		if patch.Name != nil {
			current.Name = *patch.Name
		}
		if patch.Location != nil {
			current.Location = *patch.Location
		}
		if patch.Date != nil {
			current.Date = *patch.Date
		}
	}

	updated, err := s.repository.UpdateOne(ctx, eventId, current)
	if err != nil {
		return nil, models.InternalError
	}
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
