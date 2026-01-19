package repositories

import (
	"context"
	"errors"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) models.EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) GetMany(ctx context.Context) ([]*models.Event, error) {
	events := []*models.Event{}

	err := r.db.WithContext(ctx).
		Table("events").
		Select("events.*, " +
			"COUNT(tickets.id) as total_tickets_purchased, " +
			"SUM(CASE WHEN tickets.entered = true THEN 1 ELSE 0 END) as total_tickets_entered").
		Joins("LEFT JOIN tickets ON tickets.event_id = events.id").
		Group("events.id").
		Scan(&events).Error

	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *EventRepository) GetOne(ctx context.Context, eventId uint) (*models.Event, error) {
	event := &models.Event{}

	err := r.db.WithContext(ctx).
		Table("events").
		Select("events.*, "+
			"COUNT(tickets.id) as total_tickets_purchased, "+
			"SUM(CASE WHEN tickets.entered = true THEN 1 ELSE 0 END) as total_tickets_entered").
		Joins("LEFT JOIN tickets ON tickets.event_id = events.id").
		Where("events.id = ?", eventId).
		Group("events.id").
		Scan(event).Error

	if err != nil {
		return nil, err
	}

	if event.ID == 0 {
		return nil, models.ErrRecordNotFound
	}

	return event, nil
}

func (r *EventRepository) CreateOne(ctx context.Context, event *models.Event) (*models.Event, error) {
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		return nil, err
	}
	return event, nil
}

func (r *EventRepository) UpdateOne(ctx context.Context, eventId uint, event *models.Event) (*models.Event, error) {
	existing := &models.Event{}
	if err := r.db.WithContext(ctx).First(existing, eventId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrRecordNotFound
		}
		return nil, err
	}

	event.ID = existing.ID

	if err := r.db.WithContext(ctx).Save(event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

func (r *EventRepository) DeleteOne(ctx context.Context, eventId uint) error {
	res := r.db.WithContext(ctx).Delete(&models.Event{}, eventId)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return models.ErrRecordNotFound
	}
	return nil
}
