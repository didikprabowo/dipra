package dipra

import (
	"fmt"
	"net/http"
	"time"
)

// WrapError for handler error
type WrapError struct {
	Code     int         `json:"code"`
	Message  interface{} `json:"message"`
	Internal string      `json:"-"`
	Date     time.Time   `json:"-"`
}

var (
	// Err404 for hander 404
	Err404 *WrapError = &WrapError{
		Code:    http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	}

	// Err405 for hander 405
	Err405 *WrapError = &WrapError{
		Code:    http.StatusMethodNotAllowed,
		Message: http.StatusText(http.StatusMethodNotAllowed),
	}

	// Err500 for hander 500
	Err500 *WrapError = &WrapError{
		Code:    http.StatusMethodNotAllowed,
		Message: http.StatusText(http.StatusInternalServerError),
	}
)

func (e *WrapError) Error() string {
	return fmt.Sprintf("Code : %d, detail : %v", e.Code, e.Message)
}

func (e *WrapError) String() string {
	return e.Message.(string)
}
