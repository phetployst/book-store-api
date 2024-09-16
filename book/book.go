package book

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/phetployst/book-store-api/middleware"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model `json:"-" swaggerignore:"true"`
	Title      string `json:"title" validate:"required"`
	Author     string `json:"author" validate:"required"`
	ISBN       string `json:"isbn" validate:"required,isbn"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (c *CustomValidator) Validate(i interface{}) error {
	if err := c.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func validateISBN(fl validator.FieldLevel) bool {
	isbn := fl.Field().String()
	re := regexp.MustCompile(`^\d{10}(\d{3})?$`)
	return re.MatchString(isbn)
}

type handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *handler {
	return &handler{db: db}
}

// Create godoc
// @Summary Add a new book
// @Description Creates a new book and stores it in the database. The book object must pass validation before being saved.
// @Tags books
// @Accept  json
// @Produce  json
// @Param book body Book true "New book object"
// @Success 201 {object} Book "Created book"
// @Failure 400 {object} map[string]string "Validation failed or failed to bind data"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /books [post]
func (handler *handler) Create(c echo.Context) error {
	book := Book{}

	validator := validator.New()
	validator.RegisterValidation("isbn", validateISBN)

	c.Echo().Validator = &CustomValidator{validator: validator}
	logger := middleware.GetLogger(c)

	if err := c.Bind(&book); err != nil {
		logger.Error("failed to bind book", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(book); err != nil {
		logger.Error("failed to validate book", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if result := handler.db.Create(&book); result.Error != nil {
		logger.Error("failed to insert book", zap.Error(result.Error))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}

	logger.Info("book created", zap.Any("book", book))
	return c.JSON(http.StatusCreated, book)

}

// GetAll godoc
// @Summary Get all books
// @Description Fetch a list of all books from the database. This endpoint retrieves the complete list of books without any filters.
// @Tags books
// @Accept  json
// @Produce  json
// @Success 200 {array} Book "List of books"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /books [get]
func (handler *handler) GetAll(c echo.Context) error {
	var books []Book
	if result := handler.db.Find(&books); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}

	return c.JSON(http.StatusOK, books)

}

// GetById godoc
// @Summary Retrieve a book by its ID
// @Description Fetches details of a specific book by its unique ID. If the book is not found, it returns a 404 error.
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path int true "Book ID"
// @Success 200 {object} Book "Book details"
// @Failure 404 {object} map[string]string "Book not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /books/{id} [get]
func (handler *handler) GetById(c echo.Context) error {
	book := Book{}
	id := c.Param("id")

	result := handler.db.First(&book, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Book not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, book)
}

// Update godoc
// @Summary Update an existing book
// @Description Updates the details of an existing book. The book must exist, and the request body should pass validation checks.
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path int true "Book ID"
// @Param book body Book true "Updated book object"
// @Success 200 {object} Book "Updated book details"
// @Failure 400 {object} map[string]string "Validation failed or failed to bind data"
// @Failure 404 {object} map[string]string "Book not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /books/{id} [put]
func (handler *handler) Update(c echo.Context) error {
	book := Book{}
	id := c.Param("id")

	validator := validator.New()
	validator.RegisterValidation("isbn", validateISBN)

	c.Echo().Validator = &CustomValidator{validator: validator}
	logger := middleware.GetLogger(c)

	if err := handler.db.First(&book, id).Error; err != nil {
		logger.Error("book not found", zap.String("id", id), zap.Error(err))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
	}

	if err := c.Bind(&book); err != nil {
		logger.Error("failed to bind book", zap.String("id", id), zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to bind book data"})
	}

	if err := c.Validate(book); err != nil {
		logger.Error("failed to validate book", zap.Any("book", book), zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	if result := handler.db.Save(&book); result.Error != nil {
		logger.Error("failed to update book", zap.Any("book", book), zap.Error(result.Error))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update book"})
	}

	logger.Info("book updated successfully", zap.Any("book", book))
	return c.JSON(http.StatusOK, book)
}

// Delete godoc
// @Summary Delete a book by its ID
// @Description Deletes a book by its unique ID. If the book is not found, it returns a 404 error. Otherwise, it returns a success message.
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]string "Book successfully deleted"
// @Failure 404 {object} map[string]string "Book not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /books/{id} [delete]
func (handler *handler) Delete(c echo.Context) error {
	book := Book{}
	id := c.Param("id")

	result := handler.db.Delete(&book, id)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Book successfully deleted"})
}
