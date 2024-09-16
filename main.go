package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/phetployst/book-store-api/book"
	"github.com/phetployst/book-store-api/config"
	"github.com/phetployst/book-store-api/middleware"
	"github.com/phetployst/book-store-api/router"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/phetployst/book-store-api/docs"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Book Store API
// @version 1.0
// @description This is a RESTful API for managing books in a book store. It supports CRUD operations such as creating, retrieving, updating, and deleting books.

// @contact.name API Support Team
// @contact.email st.phetploy@gmail.com

// @host localhost:1323
// @BasePath /

// @schemes http https
func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(middleware.LogMiddleware(logger))

	osGetter := &config.OsEnvGetter{}

	configProvider := config.ConfigProvider{Getter: osGetter}
	config := configProvider.GetConfig()

	db, err := gorm.Open(postgres.Open(config.Server.DBConnectionString), &gorm.Config{
		Logger: middleware.CreateGormLogger(),
	})
	if err != nil {
		logger.Fatal("failed to open database connection", zap.Error(err))
		panic("failed to connect to database")
	}

	db.AutoMigrate(&book.Book{})
	router.RegisterRoutes(e, db)
	address := fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal("failed to shutdown server", zap.Error(err))
	}
}
