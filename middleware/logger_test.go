package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogMiddleware(t *testing.T) {
	t.Run("should set logger with default parent and span id to context", func(t *testing.T) {
		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		logger := zap.New(observedZapCore)
		e := echo.New()
		e.Use(LogMiddleware(logger))

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		handler := logHelloWorldHandler()
		err := log(handler, logger)(c)

		assert.NoError(t, err)
		assertLoggerInContext(t, c)
		got := getFieldsInLogContext(t, observedLogs)
		assert.NotEmpty(t, got[parentIDLogField], "expected parent-id in log context but got none")
		assert.NotEmpty(t, got[spanIDLogField], "expected span-id in log context but got none")
	})

	t.Run("should set parent id with X-Request-ID to context given X-Request-ID exists", func(t *testing.T) {
		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		logger := zap.New(observedZapCore)
		e := echo.New()
		e.Use(LogMiddleware(logger))

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.Request().Header.Set("X-Request-ID", "request-id")
		handler := logHelloWorldHandler()
		err := log(handler, logger)(c)

		assert.NoError(t, err)
		assertLoggerInContext(t, c)
		got := getFieldsInLogContext(t, observedLogs)
		assert.Equal(t, "request-id", got[parentIDLogField], "expected parent-id to be request-id but got %s", got[parentIDLogField])
		assert.NotEmpty(t, got[spanIDLogField], "expected span-id in log context but got none")
	})
}

func TestGetLogger(t *testing.T) {
	t.Run("should return defined logger given logger context is set", func(t *testing.T) {
		logger, _ := zap.NewProduction()
		e := echo.New()
		c := e.NewContext(nil, nil)
		c.Set(loggerContextKey, logger)

		got := GetLogger(c)

		assert.Equal(t, logger, got)
	})

	t.Run("should return default no-op logger given logger context is not set", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(nil, nil)

		got := GetLogger(c)
		want := zap.NewNop()

		assert.Equal(t, want, got)
	})
}

func assertLoggerInContext(t testing.TB, c echo.Context) {
	t.Helper()
	assert.NotNil(t, GetLogger(c), "expected logger in context but got none")
}

func getFieldsInLogContext(t testing.TB, observedLogs *observer.ObservedLogs) map[string]string {
	t.Helper()
	got := make(map[string]string)
	for _, f := range observedLogs.All()[0].Context {
		got[f.Key] = f.String
	}

	return got
}

func logHelloWorldHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := GetLogger(c)
		logger.Info("Hello, World!")
		return c.String(http.StatusOK, "Hello, World!")
	}
}
