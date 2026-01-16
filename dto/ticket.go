package dto

type ValidateTicket struct {
	TicketID uint `json:"ticket_id"`
}

type CreateTicketInput struct {
	EventID uint `json:"event_id" validate:"required,gt=0"`
}
