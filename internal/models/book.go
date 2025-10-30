package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Book represents a book in the library
type Book struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Title          string    `gorm:"type:varchar(255);not null;index" json:"title" validate:"required"`
	ISBN           string    `gorm:"type:varchar(13);unique;not null;index" json:"isbn" validate:"required,len=13"`
	AuthorID       uuid.UUID `gorm:"type:uuid;not null" json:"author_id"`
	Author         Author    `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"author"`
	PublisherID    uuid.UUID `gorm:"type:uuid;not null" json:"publisher_id"`
	Publisher      Publisher `gorm:"foreignKey:PublisherID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"publisher"`
	Year           int       `gorm:"not null" json:"year" validate:"required,min=1000,max=9999"`
	Genre          string    `gorm:"type:varchar(100);not null;index" json:"genre" validate:"required"`
	Quantity       int       `gorm:"not null" json:"quantity" validate:"required,min=0"`
	QuantityIssued int       `gorm:"not null;default:0" json:"quantity_issued" validate:"min=0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (b *Book) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}
