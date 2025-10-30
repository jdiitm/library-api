package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/library-api/internal/api"
	"github.com/library-api/internal/models"
	"github.com/library-api/internal/repository"
	"github.com/library-api/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var router *gin.Engine

func TestMain(m *testing.M) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=library_test port=5434 sslmode=disable"
	}
	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to test DB:", err)
		os.Exit(1)
	}

	// Migrate the schema
	err = testDB.AutoMigrate(&models.Author{}, &models.Publisher{}, &models.Book{})
	if err != nil {
		fmt.Println("Failed to migrate:", err)
		os.Exit(1)
	}

	// Set up repository, service, handler
	bookRepo := repository.NewBookRepository(testDB)
	bookService := service.NewBookService(bookRepo)
	bookHandler := api.NewBookHandler(bookService)

	router = gin.Default()
	router.GET("/books", bookHandler.ListBooks)
	router.POST("/books", bookHandler.CreateBook)
	router.GET("/books/:id", bookHandler.GetBook)
	router.PUT("/books/:id", bookHandler.UpdateBook)
	router.DELETE("/books/:id", bookHandler.DeleteBook)
	router.POST("/books/:id/issue", bookHandler.IssueBook)
	router.POST("/books/:id/return", bookHandler.ReturnBook)

	code := m.Run()
	os.Exit(code)
}

// Helper to clear tables before test
func clearTables() {
	testDB.Exec("TRUNCATE TABLE books RESTART IDENTITY CASCADE;")
	testDB.Exec("TRUNCATE TABLE authors RESTART IDENTITY CASCADE;")
	testDB.Exec("TRUNCATE TABLE publishers RESTART IDENTITY CASCADE;")
}

// helper to create unique ISBN
func uniqueISBN() string {
	return fmt.Sprintf("%013d", time.Now().UnixNano()%1_000_000_000_0000)
}

// Helper to create a test book with valid author and publisher
func createTestBook(t *testing.T) models.Book {
	author := models.Author{Name: "Integration Author"}
	publisher := models.Publisher{Name: "Integration Publisher"}

	err := testDB.Create(&author).Error
	assert.NoError(t, err)
	err = testDB.Create(&publisher).Error
	assert.NoError(t, err)

	book := models.Book{
		Title:       "Integration Test Book",
		ISBN:        uniqueISBN(),
		AuthorID:    author.ID,
		PublisherID: publisher.ID,
		Year:        2025,
		Genre:       "Fiction",
		Quantity:    2,
	}

	err = testDB.Create(&book).Error
	assert.NoError(t, err)
	return book
}

func TestBookCRUDIntegration(t *testing.T) {
	clearTables()
	// Create author and publisher first
	author := models.Author{Name: "CRUD Author"}
	err := testDB.Create(&author).Error
	assert.NoError(t, err)

	publisher := models.Publisher{Name: "CRUD Publisher"}
	err = testDB.Create(&publisher).Error
	assert.NoError(t, err)

	// Create book
	reqBody := api.CreateBookRequest{
		Title:       "New Integration Book",
		ISBN:        uniqueISBN(),
		AuthorID:    author.ID,
		PublisherID: publisher.ID,
		Year:        2025,
		Genre:       "Fiction",
		Quantity:    5,
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var created models.Book
	err = json.Unmarshal(w.Body.Bytes(), &created)
	assert.NoError(t, err)
	assert.Equal(t, reqBody.Title, created.Title)

	// Get
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/books/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Update
	updateReq := api.UpdateBookRequest{
		Title:       "Updated Title",
		ISBN:        created.ISBN,
		AuthorID:    created.AuthorID,
		PublisherID: created.PublisherID,
		Year:        created.Year,
		Genre:       created.Genre,
		Quantity:    created.Quantity,
	}

	body, _ = json.Marshal(updateReq)
	fmt.Println(created)
	fmt.Println(updateReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/books/"+created.ID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var updated models.Book
	err = json.Unmarshal(w.Body.Bytes(), &updated)
	assert.NoError(t, err)
	assert.Equal(t, updateReq.Title, updated.Title)

	// Delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/books/"+updated.ID.String(), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}

func TestIssueAndReturnIntegration(t *testing.T) {
	book := createTestBook(t)
	assert.Equal(t, 0, book.QuantityIssued)

	// Issue
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books/"+book.ID.String()+"/issue", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var issued models.Book
	err := json.Unmarshal(w.Body.Bytes(), &issued)
	assert.NoError(t, err)
	assert.Equal(t, 1, issued.QuantityIssued)

	// Return
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/books/"+book.ID.String()+"/return", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var returned models.Book
	err = json.Unmarshal(w.Body.Bytes(), &returned)
	assert.NoError(t, err)
	assert.Equal(t, 0, returned.QuantityIssued)
}
