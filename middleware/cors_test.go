package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didikprabowo/dipra"
	. "github.com/didikprabowo/dipra/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCorsOrigin(t *testing.T) {
	d := dipra.Default()
	t.Run("allow-origin", func(t *testing.T) {
		t.Run("allow-*", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()

			h := CORS()(func(c *dipra.Context) error { return dipra.Err404 })

			req.Header.Set(string(dipra.AccessControllOrigin), "localhost")
			req.Header.Set(string(dipra.AccessControllOrigin), "https://kodingin.com")

			ctx := d.InitialContext(resp, req)
			h(ctx)
			assert.Equal(t, "*", resp.Header().Get(string(dipra.AccessControllOrigin)))
		})

		t.Run("allow-localhost", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()

			h := CorsWithConifg(CORSConfig{
				AllowOrigins: []string{"localhost"},
			})(func(c *dipra.Context) error { return dipra.Err404 })

			req.Header.Set(string(dipra.AccessControllOrigin), "localhost")

			ctx := d.InitialContext(resp, req)
			h(ctx)
			assert.Equal(t, "localhost", resp.Header().Get(string(dipra.AccessControllOrigin)))
		})
	})
}
