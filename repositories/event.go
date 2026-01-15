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

	res := r.db.WithContext(ctx).Model(&models.Event{}).Find(&events)
	if res.Error != nil {
		return nil, res.Error
	}
	return events, nil
}

func (r *EventRepository) GetOne(ctx context.Context, eventId uint) (*models.Event, error) {
	event := &models.Event{}

	err := r.db.WithContext(ctx).
		First(event, eventId).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
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
		return gorm.ErrRecordNotFound
	}
	return nil
}
