package book

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const createBookQuery = `INSERT INTO "books" ("created_at","updated_at","deleted_at","title","author","isbn") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`

func TestCreateBook(t *testing.T) {
	t.Run("create book given valid book", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		body := `{"title": "Designing Your Life", "author": "Bill Burnett and Dave Evans", "isbn": "9781101875322"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

		mock.ExpectBegin()
		row := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(createBookQuery).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Designing Your Life", "Bill Burnett and Dave Evans", "9781101875322").
			WillReturnRows(row)
		mock.ExpectCommit()

		handler := NewHandler(gormDB)
		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.Code)

	})

}
