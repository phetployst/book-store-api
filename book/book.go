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
	gorm.Model
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
	ISBN   string `json:"isbn" validate:"required,isbn"`
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

func (handler *handler) GetAll(c echo.Context) error {
	var books []Book
	if result := handler.db.Find(&books); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}

	return c.JSON(http.StatusOK, books)

}

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
