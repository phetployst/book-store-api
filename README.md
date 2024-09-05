# 📚 Book Store API

An API for managing a bookstore, built with Go and the Echo framework. This project allows users to browse, create, update, and delete book records, while providing efficient handling of ISBN validation and error handling.

This project is based on the [example-go-api](https://github.com/raksit31667/example-go-api) repository. Many concepts and structures were adapted and built upon from this source.

## 📝 Features

- Add, update, and delete books
- ISBN-10 and ISBN-13 validation
- Middleware for logging with request tracing (parent & span IDs)
- Error handling and proper logging
- Mock testing and CI/CD integration

## 🚀 Getting Started

### Prerequisites

- **Go** (version 1.22.1)
- **Docker** (optional, for containerization)
- **PostgreSQL** (or any other database of your choice)

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/phetployst/book-store-api.git
    ```

2. Navigate into the project directory:

    ```bash
    cd book-store-api
    ```

3. Install dependencies:

    ```bash
    go mod tidy
    ```

### Database Setup

Make sure you have PostgreSQL installed and running. Set the following environment variables with your database credentials:

```bash
export HOSTNAME=localhost
export PORT=1323
export DB_CONNECTION_STRING=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
```

### Running the Application
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

## 🛠️ Technologies
- Go: Backend language
- Echo: Web framework
- PostgreSQL: Database
- Docker: Containerization
- GitHub Actions: CI/CD

## 🤝 Contributing
Feel free to submit issues or pull requests! Contributions are always welcome.
