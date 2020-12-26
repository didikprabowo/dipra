package dipra_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/didikprabowo/dipra"
	"github.com/stretchr/testify/assert"
)

type (
	profile struct {
		Name string `json:"name"`
	}
)

var (
	pdata = profile{
		Name: "Didik prabowo",
	}
	pdataStr  = "{\"name\":\"Didik prabowo\"}"
	perrorStr = errors.New("Something when wrongs")
)

func TestCtx(t *testing.T) {
	d := Default()

	assert.NotNil(t, d.InitialContext(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", strings.NewReader(pdataStr))))

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(pdataStr))

	t.Run("response", func(t *testing.T) {
		t.Run("json", func(t *testing.T) {

			response := httptest.NewRecorder()
			c := d.InitialContext(response, req)

			err := c.JSON(http.StatusOK, pdata)
			if assert.NoError(t, err) {
				assert.Equal(t, http.StatusOK, response.Code)
				assert.Equal(t, string(MIMEApplicationJSON), response.Header().Get(string(HeaderContentType)))
				assert.Equal(t, pdataStr, response.Body.String())
			} else {
				t.Logf("Status Code :: Except %d, Actual %d", http.StatusOK, response.Code)
			}
		})

		t.Run("string", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			c := d.InitialContext(response, req)

			err := c.String(http.StatusOK, http.StatusText(http.StatusOK))
			if assert.NoError(t, err) {
				assert.Equal(t, http.StatusOK, response.Code)
				assert.Equal(t, string(MIMETextPlain), response.Header().Get(string(HeaderContentType)))
				assert.Equal(t, http.StatusText(http.StatusOK), response.Body.String())
			} else {
				t.Logf("Status Code :: Except %d, Actual %d", http.StatusOK, response.Code)
			}
		})

		t.Run("jsonp", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			c := d.InitialContext(response, req)

			err := c.JSONP(http.StatusOK, pdata)
			if assert.NoError(t, err) {
				assert.Equal(t, http.StatusOK, response.Code)
				assert.Equal(t, string(MIMEApplicationJavaScript), response.Header().Get(string(HeaderContentType)))
				assert.Equal(t, pdataStr, response.Body.String())
			} else {
				t.Logf("Status Code :: Except %d, Actual %d", http.StatusOK, response.Code)
			}
		})
	})

}
