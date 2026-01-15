package models

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Ticket struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	EventID   uint      `json:"event_id"`
	Event     Event     `json:"event" gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Entered   bool      `json:"entered" default:"false"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TicketRepository interface {
	GetMany(ctx context.Context) ([]*Ticket, error)
	GetOne(ctx context.Context, ticketID uint) (*Ticket, error)
	CreateOne(ctx context.Context, ticket *Ticket) (*Ticket, error)
	UpdateOne(ctx context.Context, ticketID uint, data *UpdateTicketInput) (*Ticket, error)
}

type TicketService interface {
	GetMany(ctx context.Context) ([]*Ticket, error)
	GetOne(ctx context.Context, ticketID uint) (*Ticket, error)
	CreateOne(ctx context.Context, ticket *Ticket) (*Ticket, error)
	UpdateOne(ctx context.Context, ticketID uint, input *UpdateTicketInput) (*Ticket, error)
	ValidateEntry(ctx context.Context, ticketID uint) (*Ticket, error)
}

type ValidateTicket struct {
	TicketId uint `json:"ticket_id"`
}

type UpdateTicketInput struct {
	Price      *float64 `json:"price"`
	SeatNumber *string  `json:"seat_number"`
	Status     *string  `json:"status"`
	Entered    *bool    `json:"entered"`
}

func NewValidationError(msg string) error {
	return fmt.Errorf("%w: %s", ErrValidation, msg)
}

var (
	InternalError = errors.New("internal server error")
	ErrValidation = errors.New("validation error")
)
