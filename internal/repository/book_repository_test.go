package repository_test

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/library-api/internal/models"
	"github.com/library-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=library_test port=5434 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Drop tables in reverse order to handle foreign key constraints
	if err := db.Migrator().DropTable(&models.Book{}, &models.Author{}, &models.Publisher{}); err != nil {
		t.Fatalf("Failed to drop tables: %v", err)
	}

	err = db.AutoMigrate(&models.Author{}, &models.Publisher{}, &models.Book{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestBookRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBookRepository(db)

	// Create test author and publisher
	author := &models.Author{
		Name: "Test Author",
	}
	err := db.Create(author).Error
	assert.NoError(t, err, "Failed to create test author")

	publisher := &models.Publisher{
		Name:     "Test Publisher",
		Location: "Test Location",
	}
	err = db.Create(publisher).Error
	assert.NoError(t, err, "Failed to create test publisher")

	// Test Create
	book := &models.Book{
		Title:       "Test Book",
		ISBN:        "1234567890123",
		AuthorID:    author.ID,
		PublisherID: publisher.ID,
		Year:        2025,
		Genre:       "Test",
		Quantity:    2,
	}

	err = repo.Create(book)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, book.ID)

	// Test GetByID
	found, err := repo.GetByID(book.ID)
	assert.NoError(t, err)
	assert.Equal(t, book.Title, found.Title)
	assert.Equal(t, book.ISBN, found.ISBN)

	// Test Update
	book.Title = "Updated Test Book"
	err = repo.Update(book)
	assert.NoError(t, err)

	updated, err := repo.GetByID(book.ID)
	assert.NoError(t, err)
	assert.Equal(t, book.Title, updated.Title)

	// Test Issue
	_, err = repo.IssueBook(book.ID)
	assert.NoError(t, err)
	_, err = repo.IssueBook(book.ID)
	assert.NoError(t, err)
	_, err = repo.IssueBook(book.ID)
	assert.Error(t, err)

	// Test Return
	_, err = repo.ReturnBook(book.ID)
	assert.NoError(t, err)
	_, err = repo.ReturnBook(book.ID)
	assert.NoError(t, err)
	_, err = repo.ReturnBook(book.ID)
	assert.Error(t, err)

	// Test List
	books, total, err := repo.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, len(books))

	// Test Delete
	err = repo.Delete(book.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(book.ID)
	assert.Error(t, err)

	// Clean up author and publisher
	err = db.Delete(publisher).Error
	assert.NoError(t, err)

	err = db.Delete(author).Error
	assert.NoError(t, err)
}
