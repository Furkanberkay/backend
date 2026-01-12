package repositories

import (
	"context"
	"time"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
)

type EventRepository struct {
	db any
}

func NewEventRepository(db any) models.EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) GetMany(ctx context.Context) ([]*models.Event, error) {
	events := []*models.Event{}

	events = append(events, &models.Event{
		ID:        "fwlrfne",
		Name:      "berkay",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Location:  "atasehir",
	})

	return events, nil
}

func (r *EventRepository) GetOne(ctx context.Context, userId string) (*models.Event, error) {
	return nil, nil
}
func (r *EventRepository) CreateOne(ctx context.Context, event models.Event) (*models.Event, error) {
	return nil, nil
}
