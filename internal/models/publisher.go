package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Publisher represents a book publisher
type Publisher struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;unique;index" json:"name" validate:"required"`
	Location  string    `gorm:"type:varchar(255)" json:"location"`
	Books     []Book    `gorm:"foreignKey:PublisherID" json:"books,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (p *Publisher) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}
