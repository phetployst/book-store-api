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
	"github.com/phetployst/book-store-api/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()

	osGetter := &config.OsEnvGetter{}

	configProvider := config.ConfigProvider{Getter: osGetter}
	config := configProvider.GetConfig()

	db, err := gorm.Open(postgres.Open(config.Server.DBConnectionString), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&book.Book{})
	router.RegisterRoutes(e, db)
	address := fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
