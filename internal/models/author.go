package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Author represents a book author
type Author struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;index" json:"name" validate:"required"`
	Biography string    `gorm:"type:text" json:"biography"`
	Books     []Book    `gorm:"foreignKey:AuthorID" json:"books,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (a *Author) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}
