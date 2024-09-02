package main

import (
	"fmt"

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

	e.Logger.Fatal(e.Start(address))
}
