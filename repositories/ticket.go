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

func (r *TicketRepository) GetMany(ctx context.Context, userId uint) ([]*models.Ticket, error) {
	var tickets []*models.Ticket

	err := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Preload("Event").
		Where("user_id = ?", userId).
		Order("updated_at DESC").
		Find(&tickets).Error

	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepository) GetOne(ctx context.Context, userId uint, ticketId uint) (*models.Ticket, error) {
	ticket := new(models.Ticket)

	err := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Preload("Event").
		Where("id = ? AND user_id = ?", ticketId, userId).
		First(ticket).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrRecordNotFound
		}
		return nil, err
	}

	return ticket, nil
}

func (r *TicketRepository) CreateOne(ctx context.Context, userId uint, ticket *models.Ticket) (*models.Ticket, error) {
	ticket.UserID = userId

	err := r.db.WithContext(ctx).Create(ticket).Error
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *TicketRepository) UpdateOne(ctx context.Context, userId uint, ticketId uint, data *models.UpdateTicketInput) (*models.Ticket, error) {
	var ticket models.Ticket

	result := r.db.WithContext(ctx).
		Model(&ticket).
		Clauses(clause.Returning{}).
		Where("id = ? AND user_id = ?", ticketId, userId).
		Updates(data)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, models.ErrRecordNotFound
	}

	return &ticket, nil
}
