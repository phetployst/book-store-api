# 📚 Book Store API

An API for managing a bookstore, built with Go and the Echo framework. This project allows users to browse, create, update, and delete book records.

This project is based on the [example-go-api](https://github.com/raksit31667/example-go-api) repository. Many concepts and structures were adapted and built upon from this source.

## 📦️ Packages

```bash
go get github.com/DATA-DOG/go-sqlmock
go get github.com/go-playground/validator/v10
go get github.com/google/uuid
go get github.com/labstack/echo/v4
go get github.com/stretchr/testify
go get github.com/swaggo/echo-swagger
go get github.com/swaggo/swag
go get go.uber.org/zap
go get gorm.io/driver/postgres
go get gorm.io/gorm
```

## 🚀 Running the Application
1. To start the application locally:

```bash
go run main.go
```

2. The API will be running at http://localhost:1323.

### Running Tests
To run the test suite:

```bash
go test ./...
```

## 📚 API Documentation

| Method | Endpoint        | Description          |
|--------|-----------------|----------------------|
| GET    | /books          | Get all books        |
| GET    | /books/:id      | Get a specific book  |
| POST   | /books          | Add a new book       |
| PUT    | /books/:id      | Update a book        |
| DELETE | /books/:id      | Delete a book        |

### Sample Request
To add a new book:<br>
POST /books<br>
Content-Type: application/json

```bash
{
    "title": "Clean Code",
    "author": "Robert C. Martin",
    "isbn": "9780132350884",
}
```
