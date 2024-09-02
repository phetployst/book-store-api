package router

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
)

type Route struct {
	Path   string
	Method string
}

func TestRegisterRoutes(t *testing.T) {
	e := echo.New()
	defer e.Close()

	RegisterRoutes(e, nil)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	e.ServeHTTP(response, request)

	routes := e.Routes()

	got := make([]Route, len(routes))
	for i, route := range routes {
		got[i] = Route{
			Path:   route.Path,
			Method: route.Method,
		}
	}

	want := []Route{
		{"/", http.MethodGet},
		{"/books", http.MethodPost},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v routes but want %v routes", got, want)
	}
}
