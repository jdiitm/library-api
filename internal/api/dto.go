package api

import "github.com/google/uuid"

type CreateBookRequest struct {
	Title       string    `json:"title" validate:"required"`
	ISBN        string    `json:"isbn" validate:"required,len=13"`
	AuthorID    uuid.UUID `json:"author_id" validate:"required"`
	PublisherID uuid.UUID `json:"publisher_id" validate:"required"`
	Year        int       `json:"year" validate:"required,min=1000,max=9999"`
	Genre       string    `json:"genre" validate:"required"`
	Quantity    int       `json:"quantity" validate:"required,min=0"`
}

type UpdateBookRequest struct {
	Title       string    `json:"title" validate:"required"`
	ISBN        string    `json:"isbn" validate:"required,len=13"`
	AuthorID    uuid.UUID `json:"author_id" validate:"required"`
	PublisherID uuid.UUID `json:"publisher_id" validate:"required"`
	Year        int       `json:"year" validate:"required,min=1000,max=9999"`
	Genre       string    `json:"genre" validate:"required"`
	Quantity    int       `json:"quantity" validate:"required,min=0"`
}
