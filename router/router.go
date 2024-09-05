package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/phetployst/book-store-api/book"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	bookHandler := book.NewHandler(db)

	e.POST("/books", bookHandler.Create)
	e.GET("/books", bookHandler.GetAll)
}
