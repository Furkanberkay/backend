package repositories

import (
	"context"
	"errors"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) models.TicketRepository {
	return &TicketRepository{
		db: db,
	}
}

func (r *TicketRepository) GetMany(ctx context.Context) ([]*models.Ticket, error) {
	var tickets []*models.Ticket
	err := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Preload("Event").
		Order("updated_at DESC").
		Find(&tickets).Error

	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepository) GetOne(ctx context.Context, ticketId uint) (*models.Ticket, error) {
	ticket := new(models.Ticket)
	err := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Preload("Event").
		Where("id = ?", ticketId).
		First(ticket).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrRecordNotFound
		}
		return nil, err
	}

	return ticket, nil
}

func (r *TicketRepository) CreateOne(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	err := r.db.WithContext(ctx).Create(ticket).Error
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *TicketRepository) UpdateOne(ctx context.Context, ticketId uint, data *models.UpdateTicketInput) (*models.Ticket, error) {
	var ticket models.Ticket
	ticket.ID = ticketId

	err := r.db.WithContext(ctx).
		Model(&ticket).
		Clauses(clause.Returning{}).
		Updates(data).Error

	if err != nil {
		return nil, err
	}

	if ticket.ID == 0 {
		return nil, models.ErrRecordNotFound
	}

	return &ticket, nil
}
