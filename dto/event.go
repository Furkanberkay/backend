package dto

import "time"

type EventRequest struct {
	ID       uint      `json:"id" gorm:"id"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Date     time.Time `json:"date"`
}
