package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/library-api/internal/api"
	"github.com/library-api/internal/models"
	"github.com/library-api/internal/repository"
	"github.com/library-api/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate database tables
	err = db.AutoMigrate(&models.Author{}, &models.Publisher{}, &models.Book{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Add a CHECK constraint for quantity_issued validity
	if !db.Migrator().HasConstraint(&models.Book{}, "chk_quantity_issued_valid") {
		err = db.Exec(`
		ALTER TABLE books 
		ADD CONSTRAINT chk_quantity_issued_valid 
		CHECK (quantity_issued >= 0 AND quantity_issued <= quantity);
	`).Error
		if err != nil {
			log.Fatalf("Failed to create constraint chk_quantity_issued_valid: %v", err)
		}
	}

	// Initialize repositories
	bookRepo := repository.NewBookRepository(db)

	// Initialize services
	bookService := service.NewBookService(bookRepo)

	// Initialize handlers
	bookHandler := api.NewBookHandler(bookService)

	// Setup Gin router
	r := gin.Default()

	// Routes
	v1 := r.Group("/api/v1")
	{
		books := v1.Group("/books")
		{
			books.GET("", bookHandler.ListBooks)
			books.GET("/:id", bookHandler.GetBook)
			books.POST("", bookHandler.CreateBook)
			books.PUT("/:id", bookHandler.UpdateBook)
			books.DELETE("/:id", bookHandler.DeleteBook)
			books.POST("/:id/issue", bookHandler.IssueBook)
			books.POST("/:id/return", bookHandler.ReturnBook)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
