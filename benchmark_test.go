package dipra

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkSingleRoute ...
func BenchmarkSingleRoute(b *testing.B) {
	e := Default()
	e.GET("/", func(c *Context) error {
		return c.String(http.StatusOK, "OK")
	})

	RunningRequest(b, http.MethodGet, "/", e)
}

// BenchmarkCors ...
func BenchmarkCors(b *testing.B) {
	e := Default()
	e.Use(CORS())
	e.GET("/", func(c *Context) error {
		return nil
	})
	RunningRequest(b, http.MethodGet, "/", e)
}

func BenchmarkMiddlewareHandler(b *testing.B) {
	e := Default()

	e.GET("/", func(c *Context) error {
		return nil
	}, func(next HandlerFunc) HandlerFunc {
		return func(cs *Context) error {
			return next(cs)
		}
	})

	RunningRequest(b, http.MethodGet, "/", e)
}

// BenchmarkLoggerMiddleware ...
func BenchmarkLoggerMiddleware(b *testing.B) {
	e := Default()
	e.Use(Logger())
	e.GET("/", func(c *Context) error {
		return nil
	})
	RunningRequest(b, http.MethodGet, "/", e)
}

// RunningRequest ...
func RunningRequest(b *testing.B, method string, path string, e *Engine) {

	request := httptest.NewRequest(method, path, nil)
	res := httptest.NewRecorder()
	e.ServeHTTP(res, request)

	for i := 0; i < b.N; i++ {
		e.ServeHTTP(res, request)
	}
}
