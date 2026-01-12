package services

import (
	"context"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
)

type EventService struct {
	repository models.EventRepository
}

func NewEventService(repository models.EventRepository) models.EventService {
	return &EventService{
		repository: repository,
	}
}

func (s *EventService) GetMany(ctx context.Context) ([]*models.Event, error) {
	return nil, nil
}
func (s *EventService) GetOne(ctx context.Context, eventId string) (*models.Event, error) {
	return nil, nil

}
func (s *EventService) CreateOne(ctx context.Context, event models.Event) (*models.Event, error) {
	return nil, nil

}
