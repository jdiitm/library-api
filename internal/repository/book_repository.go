package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/library-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookRepository interface {
	Create(book *models.Book) error
	Update(book *models.Book) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*models.Book, error)
	List(page, limit int) ([]models.Book, int64, error)
	IssueBook(id uuid.UUID) (*models.Book, error)
	ReturnBook(id uuid.UUID) (*models.Book, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) Update(book *models.Book) error {
	// Only update scalar fields to avoid overwriting Author/Publisher relations
	return r.db.Model(&models.Book{}).
		Where("id = ?", book.ID).
		Updates(map[string]interface{}{
			"title":           book.Title,
			"isbn":            book.ISBN,
			"author_id":       book.AuthorID,
			"publisher_id":    book.PublisherID,
			"year":            book.Year,
			"genre":           book.Genre,
			"quantity":        book.Quantity,
			"quantity_issued": book.QuantityIssued,
			"updated_at":      book.UpdatedAt,
		}).Error
}

func (r *bookRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Book{}, id).Error
}

func (r *bookRepository) GetByID(id uuid.UUID) (*models.Book, error) {
	var book models.Book
	err := r.db.Preload("Author").Preload("Publisher").First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) List(page, limit int) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.Book{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Author").Preload("Publisher").
		Offset(offset).
		Limit(limit).
		Find(&books).Error
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (r *bookRepository) IssueBook(id uuid.UUID) (*models.Book, error) {
	var book models.Book

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Lock the row FOR UPDATE to prevent concurrent modifications
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&book, "id = ?", id).Error; err != nil {
			return err
		}

		// Business logic: ensure we don't issue beyond total quantity
		if book.QuantityIssued >= book.Quantity {
			return fmt.Errorf("no available copies to issue")
		}

		book.QuantityIssued++

		// Save updated state
		if err := tx.Save(&book).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *bookRepository) ReturnBook(id uuid.UUID) (*models.Book, error) {
	var book models.Book

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Lock the row FOR UPDATE
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&book, "id = ?", id).Error; err != nil {
			return err
		}

		// Ensure you can't return below 0
		if book.QuantityIssued <= 0 {
			return fmt.Errorf("no issued copies to return")
		}

		book.QuantityIssued--

		// Save updated state
		if err := tx.Save(&book).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &book, nil
}
