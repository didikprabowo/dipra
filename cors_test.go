package dipra_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/didikprabowo/dipra"
	"github.com/stretchr/testify/assert"
)

func TestCorsOrigin(t *testing.T) {
	d := Default()
	t.Run("allow-origin", func(t *testing.T) {
		t.Run("allow-*", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()

			h := CORS()(func(c *Context) error { return Err404 })

			req.Header.Set(string(AccessControllOrigin), "localhost")
			req.Header.Set(string(AccessControllOrigin), "https://kodingin.com")

			ctx := d.InitialContext(resp, req)
			h(ctx)
			assert.Equal(t, "*", resp.Header().Get(string(AccessControllOrigin)))
		})

		t.Run("allow-localhost", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()

			h := CorsWithConifg(CORSConfig{
				AllowOrigins: []string{"localhost"},
			})(func(c *Context) error { return Err404 })

			req.Header.Set(string(AccessControllOrigin), "localhost")

			ctx := d.InitialContext(resp, req)
			h(ctx)
			assert.Equal(t, "localhost", resp.Header().Get(string(AccessControllOrigin)))
		})
	})
}
