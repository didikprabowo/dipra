package dipra_test

import (
	"net/http/httptest"
	"testing"

	. "github.com/didikprabowo/dipra"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {

	d := Default()

	t.Run("troute", func(t *testing.T) {
		t.Run("single", func(t *testing.T) {
			testTableRouteData := []Route{
				Route{
					Path:   "/index",
					Method: "GET",
					Handler: func(c *Context) error {
						return c.String(200, "OK")
					},
				},
			}

			for _, v := range testTableRouteData {
				d.AddRoute(v.Path, v.Method, v.Handler)
			}

			req := httptest.NewRequest("GET", "/index", nil)
			resp := httptest.NewRecorder()
			d.ServeHTTP(resp, req)
			assert.Equal(t, "OK", resp.Body.String())
		})

		t.Run("double", func(t *testing.T) {
			testTableRouteData := []Route{
				Route{
					Path:   "/index",
					Method: "GET",
					Handler: func(c *Context) error {
						return c.String(200, "OK")
					},
				},
				Route{
					Path:   "/index",
					Method: "GET",
					Handler: func(c *Context) error {
						return c.String(200, "OKE")
					},
				},
			}

			for _, v := range testTableRouteData {
				d.AddRoute(v.Path, v.Method, v.Handler)
			}

			req := httptest.NewRequest("GET", "/index", nil)
			resp := httptest.NewRecorder()
			d.ServeHTTP(resp, req)
			assert.Equal(t, "OKE", resp.Body.String())
		})
	})

}
