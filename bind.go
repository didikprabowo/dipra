package dipra

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"

	"gopkg.in/yaml.v3"
)

type (
	// Binding ...
	Binding struct {
		// Request
		Request *http.Request
		// ContentType
		ContentType string

		// Field
		Mapper Mapper
	}
)

var (
	// ErrMimeNotSupported for handler error not content suport
	ErrMimeNotSupported = errors.New("Content-Type header is not suport")
	// ErrNotEmpty for handler mime not found
	ErrNotEmpty = errors.New("Request body must not be empty")
	// ErrBadRequest for handler panic error
	ErrBadRequest = errors.New(http.StatusText(http.StatusBadRequest))
)

// SetBind ...
func (b *Binding) SetBind(r *http.Request) {
	b.Request = r
	b.ContentType = r.Header.Get(string(HeaderContentType))
}

// ShouldJSON for wrap request json
func (b *Binding) ShouldJSON(v interface{}) (err error) {
	if b.ContentType != string(MIMEApplicationJSON) {
		return ErrMimeNotSupported
	}

	if b.Request.Body == http.NoBody {
		return ErrNotEmpty
	}

	err = json.NewDecoder(b.Request.Body).Decode(&v)
	if err != nil {
		return ErrBadRequest
	}
	return err
}

// ShouldXML for wrap request xml
func (b *Binding) ShouldXML(v interface{}) (err error) {
	if b.ContentType != string(MIMEApplicationXML) {
		return ErrMimeNotSupported
	}

	if b.Request.Body == http.NoBody {
		return ErrNotEmpty
	}

	err = xml.NewDecoder(b.Request.Body).Decode(&v)
	if err != nil {
		return ErrBadRequest
	}
	return err
}

// ShouldYAML wrap request yaml
func (b *Binding) ShouldYAML(v interface{}) (err error) {
	if !(b.ContentType == string(MIMEApplicationYAML) || b.ContentType == string(MIMETextYAML)) {
		return ErrMimeNotSupported
	}

	if b.Request.Body == http.NoBody {
		return ErrNotEmpty
	}

	err = yaml.NewDecoder(b.Request.Body).Decode(&v)
	if err != nil {
		return ErrBadRequest
	}
	return err
}

// ShouldQuery for use wrap request Query
func (b *Binding) ShouldQuery(value interface{}) (err error) {
	q := b.Request.URL.Query()
	bimap := map[string]string{}

	for k, v := range q {
		bimap[k] = v[0]
	}
	err = b.Mapper.Set(value)
	if err != nil {
		return ErrBadRequest
	}
	err = b.Mapper.MapToStruct(bimap)
	if err != nil {
		return ErrBadRequest
	}
	return err
}
