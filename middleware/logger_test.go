package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didikprabowo/dipra"
	. "github.com/didikprabowo/dipra/middleware"
)

func TestLogger(t *testing.T) {
	d := dipra.Default()

	{
		req := httptest.NewRequest("GET", "/", nil)
		resp := httptest.NewRecorder()

		c := d.InitialContext(resp, req)
		l := Logger()(func(c *dipra.Context) error {
			return c.JSON(http.StatusOK, "OK")
		})
		l(c)
	}

	{
		req := httptest.NewRequest("GET", "/", nil)
		resp := httptest.NewRecorder()

		c := d.InitialContext(resp, req)
		l := Logger()(func(c *dipra.Context) error {
			return errors.New("Something when wrongs")
		})
		l(c)
	}
}
