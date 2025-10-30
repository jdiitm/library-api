package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/library-api/internal/models"
	"github.com/library-api/internal/service"
)

type BookHandler struct {
	bookService service.BookService
	validate    *validator.Validate
}

func NewBookHandler(bookService service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
		validate:    validator.New(),
	}
}

// ListBooks godoc
// @Summary List all books
// @Description get books
// @Tags books
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} []models.Book
// @Router /books [get]
func (h *BookHandler) ListBooks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	books, total, err := h.bookService.ListBooks(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  books,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetBook godoc
// @Summary Get a book
// @Description get book by ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	book, err := h.bookService.GetBook(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// CreateBook godoc
// @Summary Create a book
// @Description create new book
// @Tags books
// @Accept  json
// @Produce  json
// @Param book body CreateBookRequest true "Create book"
// @Success 201 {object} models.Book
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := models.Book{
		Title:       req.Title,
		ISBN:        req.ISBN,
		AuthorID:    req.AuthorID,
		PublisherID: req.PublisherID,
		Year:        req.Year,
		Genre:       req.Genre,
		Quantity:    req.Quantity,
	}

	if err := h.bookService.CreateBook(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// UpdateBook godoc
// @Summary Update a book
// @Description update book by ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Param book body UpdateBookRequest true "Update book"
// @Success 200 {object} models.Book
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the existing book from DB first
	book, err := h.bookService.GetBook(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Update allowed fields
	book.Title = req.Title
	book.ISBN = req.ISBN
	book.AuthorID = req.AuthorID
	book.PublisherID = req.PublisherID
	book.Year = req.Year
	book.Genre = req.Genre
	book.Quantity = req.Quantity

	if err := h.bookService.UpdateBook(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description delete book by ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Success 204 "No Content"
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := h.bookService.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// IssueBook godoc
// @Summary Issue a book
// @Description increment the QuantityIssued of a book
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /books/{id}/issue [post]
func (h *BookHandler) IssueBook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	book, err := h.bookService.IssueBook(id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// ReturnBook godoc
// @Summary Return a book
// @Description decrement the QuantityIssued of a book
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /books/{id}/return [post]
func (h *BookHandler) ReturnBook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	book, err := h.bookService.ReturnBook(id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}
