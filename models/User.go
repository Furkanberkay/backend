package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	Manager  UserRole = "manager"
	Attendee UserRole = "attendee"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Surname   string    `gorm:"type:varchar(100);not null"`
	Birthday  time.Time `gorm:"type:date"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Role      UserRole  `gorm:"type:varchar(20);default:'attendee';not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) AfterCreate(db *gorm.DB) error {
	if u.ID == 1 {
		db.Model(&u).Update("role", Manager)
	}
	return nil
}
