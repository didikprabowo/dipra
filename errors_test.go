package dipra_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/didikprabowo/dipra"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {

	d := Default()

	err := WrapError{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}

	assert.Equal(t, `Code : 500, detail : Internal Server Error`, err.Error())

	jsonBytes, _ := json.Marshal(err)
	assert.Equal(t, `{"code":500,"message":"Internal Server Error"}`, string(jsonBytes))
	assert.Equal(t, err.String(), err.Message)
	assert.NotNil(t, jsonBytes)

	t.Run("default-error", func(t *testing.T) {
		t.Run("404", func(t *testing.T) {
			var err error = Err404
			jsonBytes404, _ := json.Marshal(err)

			r := httptest.NewRequest(http.MethodGet, "/no", nil)
			resp := httptest.NewRecorder()
			d.ServeHTTP(resp, r)
			assert.Equal(t, `{"error":`+string(jsonBytes404)+`}`, resp.Body.String())
		})

		t.Run("500", func(t *testing.T) {

			d.GET("/", func(c *Context) error {
				return Err500
			})

			jsonBytes500, _ := json.Marshal(Err500)

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()

			d.ServeHTTP(resp, r)
			assert.Equal(t, `{"error":`+string(jsonBytes500)+`}`, resp.Body.String())
		})
	})
}
