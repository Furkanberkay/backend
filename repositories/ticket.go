package repositories

import (
	"context"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"gorm.io/gorm"
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
	tickets := []*models.Ticket{}
	res := r.db.WithContext(ctx).Model(models.Ticket{}).Preload("Event").Order("updated_at DESC").Find(&tickets)
	if res.Error != nil {
		return nil, res.Error
	}
	return tickets, nil
}

func (r *TicketRepository) GetOne(ctx context.Context, ticketId uint) (*models.Ticket, error) {
	ticket := new(models.Ticket)
	resp := r.db.WithContext(ctx).Model(&models.Ticket{}).Preload("Event").Where("id=?", ticketId).First(ticket)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return ticket, nil
}

func (r *TicketRepository) CreateOne(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	if err := r.db.WithContext(ctx).Model(&models.Ticket{}).Create(ticket).Error; err != nil {
		return nil, err
	}
	return r.GetOne(ctx, ticket.ID)
}

func (r *TicketRepository) UpdateOne(ctx context.Context, ticketId uint, updateData map[string]interface{}) (*models.Ticket, error) {
	resp := r.db.WithContext(ctx).Model(&models.Ticket{}).Where("id = ?", ticketId).Updates(updateData)
	if resp.Error != nil {
		return nil, resp.Error
	}
	if resp.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return r.GetOne(ctx, ticketId)

}
