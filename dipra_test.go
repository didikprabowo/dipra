package dipra

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEngine ...
func TestEngine(t *testing.T) {
	e := Default()

	assert.NotNil(t, e)
	assert.Empty(t, e.Route)
}

// TestInitialContext ...
func TestInitialContext(t *testing.T) {
	e := Default()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	e.ServeHTTP(res, request)

	c := e.InitialContext(res, request)

	assert.Equal(t, c.GetResponse(), res)
	assert.Equal(t, c.GetRequest(), request)

	defaulterrorHTTP(res, http.StatusInternalServerError, errors.New("Interval server error"))
}

// TestMethodGet...
func TestMethodGet(t *testing.T) {
	e := Default()

	{
		e.GET("/", func(c *Context) error {
			return c.String(http.StatusOK, "OK")
		})

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		e.ServeHTTP(res, request)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "OK", string(res.Body.Bytes()))
	}
}

// TestMethodPOST ....
func TestMethodPOST(t *testing.T) {
	e := Default()
	{
		e.POST("/", func(c *Context) error {
			return c.String(http.StatusOK, "OK")
		})

		request := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()
		e.ServeHTTP(res, request)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "OK", string(res.Body.Bytes()))
	}
}

// TestMethodPUT ...
func TestMethodPUT(t *testing.T) {
	e := Default()
	{
		e.PUT("/", func(c *Context) error {
			return c.String(http.StatusOK, "OK")
		})

		request := httptest.NewRequest(http.MethodPut, "/", nil)
		res := httptest.NewRecorder()
		e.ServeHTTP(res, request)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "OK", string(res.Body.Bytes()))
	}
}

// TestMethodPATCH ...
func TestMethodPATCH(t *testing.T) {
	e := Default()
	{
		e.PATCH("/", func(c *Context) error {
			return c.String(http.StatusOK, "OK")
		})

		request := httptest.NewRequest(http.MethodPatch, "/", nil)
		res := httptest.NewRecorder()
		e.ServeHTTP(res, request)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "OK", string(res.Body.Bytes()))
	}
}

// TestMethodDELETE ...
func TestMethodDELETE(t *testing.T) {
	e := Default()
	{
		e.DELETE("/", func(c *Context) error {
			return c.String(http.StatusOK, "OK")
		})

		request := httptest.NewRequest(http.MethodDelete, "/", nil)
		res := httptest.NewRecorder()
		e.ServeHTTP(res, request)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "OK", string(res.Body.Bytes()))
	}
}

// TestMiddleware ...
func TestMiddleware(t *testing.T) {
	e := Default()
	e.Use(func(Next HandlerFunc) HandlerFunc {
		return func(c *Context) (err error) {
			return Next(c)
		}
	})

	e.GET("/", func(c *Context) error {
		return c.String(http.StatusOK, "OK")
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	e.ServeHTTP(res, request)
	assert.Equal(t, http.StatusOK, res.Code)
}
