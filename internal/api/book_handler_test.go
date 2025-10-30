package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/library-api/internal/api"
	"github.com/library-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) CreateBook(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookService) UpdateBook(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookService) DeleteBook(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookService) GetBook(id uuid.UUID) (*models.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) ListBooks(page, limit int) ([]models.Book, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookService) IssueBook(id uuid.UUID) (*models.Book, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) ReturnBook(id uuid.UUID) (*models.Book, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Book), args.Error(1)
}

func TestBookHandler_ListBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockBookService)
	handler := api.NewBookHandler(mockService)

	books := []models.Book{
		{
			Title: "Test Book 1",
			ISBN:  "1234567890123",
			Genre: "Test",
			Year:  2025,
		},
	}

	mockService.On("ListBooks", 1, 10).Return(books, int64(1), nil)

	r := gin.Default()
	r.GET("/books", handler.ListBooks)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), response["total"])
	assert.Equal(t, float64(1), response["page"])
	assert.Equal(t, float64(10), response["limit"])
}

func TestBookHandler_CreateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockBookService)
	handler := api.NewBookHandler(mockService)

	book := models.Book{
		Title:       "Test Book",
		ISBN:        "1234567890123",
		Genre:       "Fiction",
		Year:        2025,
		Quantity:    10,
		AuthorID:    uuid.New(),
		Author:      models.Author{Name: "Test Author"},
		PublisherID: uuid.New(),
		Publisher:   models.Publisher{Name: "Test Publisher"},
	}

	mockService.On("CreateBook", mock.AnythingOfType("*models.Book")).Return(nil)

	r := gin.Default()
	r.POST("/books", handler.CreateBook)

	bookJSON, _ := json.Marshal(book)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(bookJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	// Log response body for debugging
	if w.Code != 201 {
		t.Logf("Response Body: %s", w.Body.String())
	}
}

func TestBookHandler_UpdateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockBookService)
	handler := api.NewBookHandler(mockService)

	bookID := uuid.New()
	existingBook := &models.Book{
		ID:          bookID,
		Title:       "Old Title",
		ISBN:        "9876543210123",
		Genre:       "Non-Fiction",
		Year:        2024,
		Quantity:    3,
		AuthorID:    uuid.New(),
		PublisherID: uuid.New(),
	}

	// Updated data
	updateReq := api.UpdateBookRequest{
		Title:       "Updated Test Book",
		ISBN:        "1234567890123",
		Genre:       "Fiction",
		Year:        2025,
		Quantity:    5,
		AuthorID:    uuid.New(),
		PublisherID: uuid.New(),
	}

	// Mock GetBook to return the existing book
	mockService.On("GetBook", bookID).Return(existingBook, nil)
	// Mock UpdateBook to succeed
	mockService.On("UpdateBook", mock.AnythingOfType("*models.Book")).Return(nil)

	r := gin.Default()
	r.PUT("/books/:id", handler.UpdateBook)

	reqBody, _ := json.Marshal(updateReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/books/"+bookID.String(), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response models.Book
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updateReq.Title, response.Title)
	assert.Equal(t, updateReq.ISBN, response.ISBN)
	assert.Equal(t, updateReq.Genre, response.Genre)
	assert.Equal(t, updateReq.Year, response.Year)
	assert.Equal(t, updateReq.Quantity, response.Quantity)
}

func TestBookHandler_DeleteBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockBookService)
	handler := api.NewBookHandler(mockService)

	bookID := uuid.New()

	mockService.On("DeleteBook", bookID).Return(nil)

	r := gin.Default()
	r.DELETE("/books/:id", handler.DeleteBook)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/books/"+bookID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

func TestBookHandler_GetBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockBookService)
	handler := api.NewBookHandler(mockService)

	bookID := uuid.New()
	book := models.Book{
		ID:          bookID,
		Title:       "Test Book",
		ISBN:        "1234567890123",
		Genre:       "Fiction",
		Year:        2025,
		Quantity:    10,
		AuthorID:    uuid.New(),
		PublisherID: uuid.New(),
	}

	mockService.On("GetBook", bookID).Return(&book, nil)

	r := gin.Default()
	r.GET("/books/:id", handler.GetBook)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/"+bookID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response models.Book
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, book.Title, response.Title)
	assert.Equal(t, book.ISBN, response.ISBN)
}

func TestBookHandler_IssueBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockBookService)
	handler := api.NewBookHandler(mockService)

	bookID := uuid.New()
	book := &models.Book{
		ID:             bookID,
		Title:          "Test Book",
		ISBN:           "1234567890123",
		Quantity:       2,
		QuantityIssued: 1,
	}

	mockService.On("IssueBook", bookID).Return(book, nil)

	r := gin.Default()
	r.POST("/books/:id/issue", handler.IssueBook)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books/"+bookID.String()+"/issue", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response models.Book
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.QuantityIssued)
	assert.Equal(t, book.Title, response.Title)
}

func TestBookHandler_ReturnBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockBookService)
	handler := api.NewBookHandler(mockService)

	bookID := uuid.New()
	book := &models.Book{
		ID:             bookID,
		Title:          "Test Book",
		ISBN:           "1234567890123",
		Quantity:       2,
		QuantityIssued: 0,
	}

	mockService.On("ReturnBook", bookID).Return(book, nil)

	r := gin.Default()
	r.POST("/books/:id/return", handler.ReturnBook)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books/"+bookID.String()+"/return", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response models.Book
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, response.QuantityIssued)
	assert.Equal(t, book.Title, response.Title)
}
