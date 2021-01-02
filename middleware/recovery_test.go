package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didikprabowo/dipra"
	. "github.com/didikprabowo/dipra/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	d := dipra.Default()

	d.Use(Recovery())

	d.GET("/tes_panic", func(c *dipra.Context) error {
		panic("tes panic")
	})
	req := httptest.NewRequest(http.MethodGet, "/tes_panic", nil)
	w := httptest.NewRecorder()
	d.ServeHTTP(w, req)

	err := struct {
		Err dipra.WrapError `json:"error"`
	}{}

	json.Unmarshal(w.Body.Bytes(), &err)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "tes panic", err.Err.Message)
}
