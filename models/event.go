package models

import (
	"context"
	"errors"
	"time"
)

type Event struct {
	ID                    uint      `json:"id" gorm:"primaryKey"`
	Name                  string    `json:"name"`
	Location              string    `json:"location"`
	Date                  time.Time `json:"date"`
	TotalTicketsPurchased int64     `json:"totalTicketsPurchased" gorm:"-"`
	TotalTicketsEntered   int64     `json:"totalTicketsEntered" gorm:"-"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

type EventUpdateInput struct {
	Name     *string
	Location *string
	Date     *time.Time
}

type EventRepository interface {
	GetMany(ctx context.Context) ([]*Event, error)
	GetOne(ctx context.Context, eventId uint) (*Event, error)
	CreateOne(ctx context.Context, event *Event) (*Event, error)
	UpdateOne(ctx context.Context, eventId uint, event *Event) (*Event, error)
	DeleteOne(ctx context.Context, eventId uint) error
}

type EventService interface {
	GetMany(ctx context.Context) ([]*Event, error)
	GetOne(ctx context.Context, eventId uint) (*Event, error)
	CreateOne(ctx context.Context, event *Event) (*Event, error)
	UpdateOne(ctx context.Context, eventId uint, event *Event) (*Event, error)
	PatchOne(ctx context.Context, eventId uint, patch *EventUpdateInput) (*Event, error)
	DeleteOne(ctx context.Context, eventId uint) error
}

var ErrRecordNotFound = errors.New("event not found")
