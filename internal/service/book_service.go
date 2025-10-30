package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/library-api/internal/models"
	"github.com/library-api/internal/repository"
)

type BookService interface {
	CreateBook(book *models.Book) error
	UpdateBook(book *models.Book) error
	DeleteBook(id uuid.UUID) error
	GetBook(id uuid.UUID) (*models.Book, error)
	ListBooks(page, limit int) ([]models.Book, int64, error)
	IssueBook(id uuid.UUID) (*models.Book, error)
	ReturnBook(id uuid.UUID) (*models.Book, error)
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) CreateBook(book *models.Book) error {
	return s.repo.Create(book)
}

func (s *bookService) UpdateBook(book *models.Book) error {
	if book.Quantity < book.QuantityIssued {
		return fmt.Errorf("quantity cannot be less than the number of issued copies")
	}
	return s.repo.Update(book)
}

func (s *bookService) DeleteBook(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *bookService) GetBook(id uuid.UUID) (*models.Book, error) {
	return s.repo.GetByID(id)
}

func (s *bookService) ListBooks(page, limit int) ([]models.Book, int64, error) {
	return s.repo.List(page, limit)
}

func (s *bookService) IssueBook(id uuid.UUID) (*models.Book, error) {
	book, err := s.repo.IssueBook(id)
	if err != nil {
		return nil, fmt.Errorf("failed to issue book: %w", err)
	}
	return book, nil
}

func (s *bookService) ReturnBook(id uuid.UUID) (*models.Book, error) {
	book, err := s.repo.ReturnBook(id)
	if err != nil {
		return nil, fmt.Errorf("failed to return book: %w", err)
	}
	return book, nil
}
