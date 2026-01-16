package dto

import "time"

type CreateEventInput struct {
	Name     string    `json:"name" validate:"required,min=3,max=100"`
	Location string    `json:"location" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
}

type UpdateEventInput struct {
	Name     string    `json:"name" validate:"required,min=3,max=100"`
	Location string    `json:"location" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
}

type EventPatchRequest struct {
	Name     *string    `json:"name" validate:"omitempty,min=3,max=100"`
	Location *string    `json:"location" validate:"omitempty"`
	Date     *time.Time `json:"date" validate:"omitempty"`
}
