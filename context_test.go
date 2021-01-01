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
	pdataStr  = `{"name":"Didik prabowo"}`
	pDataXML  = `<profile><Name>Didik prabowo</Name></profile>`
	resperror = `{"error":{"code":500,"message":"Something when wrongs"}}`
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

		t.Run("xml", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			c := d.InitialContext(response, req)

			err := c.XML(http.StatusOK, pdata)
			if assert.NoError(t, err) {
				assert.Equal(t, http.StatusOK, response.Code)
				assert.Equal(t, string(MIMEApplicationXML), response.Header().Get(string(HeaderContentType)))
				assert.Equal(t, pDataXML, response.Body.String())
			} else {
				t.Logf("Status Code :: Except %d, Actual %d", http.StatusOK, response.Code)
			}
		})

		t.Run("error", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			c := d.InitialContext(response, req)

			c.SetError(errors.New("Something when wrongs"))
			assert.Equal(t, http.StatusInternalServerError, response.Code)
			assert.Equal(t, resperror, response.Body.String())
		})
	})
}

func TestParam(t *testing.T) {
	d := Default()
	req := httptest.NewRequest(http.MethodGet, "/user/OK/Created", strings.NewReader(pdataStr))
	response := httptest.NewRecorder()

	d.AddRoute("/user/:name/:status", http.MethodGet, nil)
	c := d.InitialContext(response, req)
	d.HandlerRoute(c)

	c.String(http.StatusOK, c.Param("name")+"AND"+c.Param("status"))
	assert.Equal(t, http.StatusText(http.StatusOK)+"AND"+http.StatusText(http.StatusCreated), response.Body.String())
	assert.Equal(t, "/user/:name/:status", c.GetPatcher())
	assert.Equal(t, "/user/OK/Created", c.GetPath())
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestQuery(t *testing.T) {
	d := Default()
	req := httptest.NewRequest(http.MethodGet, "/?name=didik", strings.NewReader(pdataStr))
	response := httptest.NewRecorder()

	c := d.InitialContext(response, req)

	c.String(http.StatusOK, c.Query("name"))
	assert.Equal(t, "didik", response.Body.String())
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestShouldJSON(t *testing.T) {
	d := Default()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(pdataStr))
	res := httptest.NewRecorder()
	req.Header.Add(string(HeaderContentType), string(MIMEApplicationJSON))
	c := d.InitialContext(res, req)
	c.SetBind(req)

	var (
		p profile
	)

	err := c.ShouldJSON(&p)
	assert.NoError(t, err)
	assert.Equal(t, pdata, p)
}
