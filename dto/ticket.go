package dto

import "time"

type CreateTicketInput struct {
	EventID uint `json:"event_id" validate:"required,gt=0"`
}

type TicketRequest struct {
	ID      uint `json:"id" gorm:"primaryKey"`
	EventID uint `json:"event_id"`
	UserID  uint `json:"userID"`
	Entered bool `json:"entered" default:"false"`
}

type TicketResponse struct {
	ID           uint      `json:"id"`
	EventID      uint      `json:"event_id"`
	EventName    string    `json:"event_name"`
	UserID       uint      `json:"user_id"`
	Quantity     int       `json:"quantity"`
	PurchaseDate time.Time `json:"purchase_date"`
	QRCode       string    `json:"qr_code"`
}
