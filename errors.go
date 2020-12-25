package dipra

import (
	"fmt"
	"time"
)

// WrapError for handler error
type WrapError struct {
	Code     int         `json:"code"`
	Message  interface{} `json:"message"`
	Internal string      `json:"-"`
	Date     time.Time   `json:"-"`
}

func (e *WrapError) Error() string {
	return fmt.Sprintf("Code : %d, detail : %v", e.Code, e.Message)
}

func (e *WrapError) String() string {
	return e.Message.(string)
}
