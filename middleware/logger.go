package middleware

import (
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

const (
	loggerContextKey = "logger"
	parentIDLogField = "parent-id"
	spanIDLogField   = "span-id"
)

func CreateGormLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)
}

func GetLogger(c echo.Context) *zap.Logger {
	switch logger := c.Get(loggerContextKey).(type) {
	case *zap.Logger:
		return logger
	default:
		return zap.NewNop()
	}
}

func LogMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return setRequestLogger(next, logger)
	}
}

func setRequestLogger(next echo.HandlerFunc, logger *zap.Logger) func(c echo.Context) error {
	return func(c echo.Context) error {
		c.Set(loggerContextKey, loggerWithParentAndSpanID(c, logger))
		return next(c)
	}
}

func loggerWithParentAndSpanID(c echo.Context, logger *zap.Logger) *zap.Logger {
	parentID := c.Request().Header.Get("X-Request-ID")
	if parentID == "" {
		parentID = uuid.New().String()
	}
	spanID := uuid.New().String()
	return logger.With(zap.String(parentIDLogField, parentID), zap.String(spanIDLogField, spanID))
}
