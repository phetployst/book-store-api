package book

import (
	"errors"
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

const (
	createBookQuery  = `INSERT INTO "books" ("created_at","updated_at","deleted_at","title","author","isbn") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
	getAllBookQuery  = `SELECT * FROM "books" WHERE "books"."deleted_at" IS NULL`
	getBookByIdQuery = `SELECT * FROM "books" WHERE "books"."id" = $1 AND "books"."deleted_at" IS NULL ORDER BY "books"."id" LIMIT $2`
	updateBookQuery  = `UPDATE "books" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"title"=$4,"author"=$5,"isbn"=$6 WHERE "books"."deleted_at" IS NULL AND "id" = $7`
)

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

	t.Run("create book given invalid book", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"title": ""}`))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		handler := NewHandler(nil)
		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)

	})

	t.Run("create book given invalid book isbn", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		body := `{"title": "The Alchemist", "author": "Paulo Coelho", "isbn": "007"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		handler := NewHandler(nil)
		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)

	})

	t.Run("create book given error during book binding", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(``))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		handler := NewHandler(nil)
		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("create book given error during query", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		body := `{"title": "The Happiness of Pursuit", "author": "Chris Guillebeau", "isbn": "9780385348876"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

		mock.ExpectBegin()
		mock.ExpectQuery(createBookQuery).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "The Happiness of Pursuit", "Chris Guillebeau", "9780385348876").
			WillReturnError(errors.New("query error"))
		mock.ExpectRollback()

		handler := NewHandler(gormDB)
		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)

	})

}

func TestGetAllBook(t *testing.T) {
	t.Run("get all books given books exist in the database", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

		rows := sqlmock.NewRows([]string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt", "title", "author", "isbn"}).
			AddRow(1, nil, nil, nil, "Four Thousand Weeks", "Oliver Burkeman", "9781785038723").
			AddRow(2, nil, nil, nil, "Atomic Habits", "James Clear", "9781847941831").
			AddRow(3, nil, nil, nil, "The Tree of a Thousand Loves", "Sukanya Kittikhun", "9786164453819")
		mock.ExpectQuery(getAllBookQuery).WillReturnRows(rows)

		handler := NewHandler(gormDB)
		err := handler.GetAll(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("get all books given error during query", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

		mock.ExpectQuery(getAllBookQuery).WillReturnError(errors.New("query error"))
		handler := NewHandler(gormDB)
		err := handler.GetAll(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestGetBookById(t *testing.T) {

	t.Run("get book by id given a book exists in the database", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetPath("/books/:id")
		c.SetParamNames("id")
		c.SetParamValues("3")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

		row := sqlmock.NewRows([]string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt", "title", "author", "isbn"})
		row.AddRow(3, nil, nil, nil, "The Tree of a Thousand Loves", "Sukanya Kittikhun", "9786164453819")
		mock.ExpectQuery(getBookByIdQuery).WithArgs("3", 1).WillReturnRows(row)

		handler := NewHandler(gormDB)
		err := handler.GetById(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("get book by id given book does not exist", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetPath("/books/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

		mock.ExpectQuery(getBookByIdQuery).WithArgs("1", 1).WillReturnError(gorm.ErrRecordNotFound)

		handler := NewHandler(gormDB)
		err := handler.GetById(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("get book by id given error during query", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetPath("/books/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

		mock.ExpectQuery(getBookByIdQuery).WithArgs("1", 1).WillReturnError(errors.New("query error"))

		handler := NewHandler(gormDB)
		err := handler.GetById(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

}

func TestUpdateBook(t *testing.T) {
	t.Run("update book successful", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		body := `{"title": "The Tree of a Thousand Loves", "author": "Sukanya Kittikhun"}`
		request := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetPath("/books/:id")
		c.SetParamNames("id")
		c.SetParamValues("29")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

		row := sqlmock.NewRows([]string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt", "title", "author", "isbn"})
		row.AddRow(1, nil, nil, nil, "The Tree of Loves", "Phetploy", "9786164453819")
		mock.ExpectQuery(getBookByIdQuery).
			WithArgs("29", 1).
			WillReturnRows(row)

		mock.ExpectBegin()
		mock.ExpectExec(updateBookQuery).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "The Tree of a Thousand Loves", "Sukanya Kittikhun", "9786164453819", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		handler := NewHandler(gormDB)
		err := handler.Update(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("update book given book does not exist", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		body := `{"title": "The Great Gatsby", "author": "F. Scott Fitzgerald", "isbn": "9780743273565"}`
		request := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetPath("/books/:id")
		c.SetParamNames("id")
		c.SetParamValues("12")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

		mock.ExpectQuery(getBookByIdQuery).
			WithArgs("12", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		handler := NewHandler(gormDB)
		err := handler.Update(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("update book given invalid input ISBN", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		body := `{"isbn": "00"}`
		request := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetPath("/books/:id")
		c.SetParamNames("id")
		c.SetParamValues("18")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

		row := sqlmock.NewRows([]string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt", "title", "author", "isbn"})
		row.AddRow(18, nil, nil, nil, "1984", "George Orwell", "9780451524935")
		mock.ExpectQuery(getBookByIdQuery).
			WithArgs("18", 1).
			WillReturnRows(row)

		handler := NewHandler(gormDB)
		err := handler.Update(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("update book given error during query execution", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		body := `{"author": "J.D. Salinger"}`
		request := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetPath("/books/:id")
		c.SetParamNames("id")
		c.SetParamValues("29")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

		row := sqlmock.NewRows([]string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt", "title", "author", "isbn"})
		row.AddRow(29, nil, nil, nil, "The Catcher in the Rye", "J.D. Saling", "9780316769488")
		mock.ExpectQuery(getBookByIdQuery).
			WithArgs("29", 1).
			WillReturnRows(row)

		mock.ExpectExec(updateBookQuery).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "The Tree of Loves", "Phetploy", "0781101875322", 1).
			WillReturnError(errors.New("query error"))
		mock.ExpectRollback()

		handler := NewHandler(gormDB)
		err := handler.Update(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

}
