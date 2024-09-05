package router

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
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
		{"/books", http.MethodPost},
		{"/books", http.MethodGet},
	}

	sort.Slice(got, func(i, j int) bool {
		return got[i].Path < got[j].Path || (got[i].Path == got[j].Path && got[i].Method < got[j].Method)
	})
	sort.Slice(want, func(i, j int) bool {
		return want[i].Path < want[j].Path || (want[i].Path == want[j].Path && want[i].Method < want[j].Method)
	})

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v routes but want %v routes", got, want)
	}
}
