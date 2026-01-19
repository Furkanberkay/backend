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
	UserID    uint      `json:"userID"`
	Event     Event     `json:"event" gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Entered   bool      `json:"entered" default:"false"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TicketRepository interface {
	GetMany(ctx context.Context, userId uint) ([]*Ticket, error)
	GetOne(ctx context.Context, userId uint, ticketID uint) (*Ticket, error)
	CreateOne(ctx context.Context, userId uint, ticket *Ticket) (*Ticket, error)
	UpdateOne(ctx context.Context, userId uint, ticketID uint, data *UpdateTicketInput) (*Ticket, error)
}

type TicketService interface {
	GetMany(ctx context.Context, userId uint) ([]*Ticket, error)
	GetOne(ctx context.Context, userId uint, ticketID uint) (*Ticket, string, error)
	CreateOne(ctx context.Context, userId uint, ticket *Ticket) (*Ticket, error)
	UpdateOne(ctx context.Context, userId uint, ticketID uint, input *UpdateTicketInput) (*Ticket, error)
	ValidateEntry(ctx context.Context, userId uint, ticketID uint) (*Ticket, error)
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

type ValidateTicket struct {
	TicketID uint `json:"ticket_id"`
	OwnerID  uint `json:"owner_id"`
}

var (
	InternalError        = errors.New("internal server error")
	ErrValidation        = errors.New("validation error")
	ErrEventNotFound     = errors.New("event not found")
	ErrTicketAlreadyUsed = errors.New("ticket already used")
)
