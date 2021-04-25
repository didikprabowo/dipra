package dipra

import (
	"net/http"
)

var (
	allowMethod = []string{
		http.MethodGet,
		http.MethodDelete,
		http.MethodPost,
		http.MethodPatch,
		http.MethodPut,
	}
)

func (e *route) allowMethods(m string) (ok bool) {
	for _, mlist := range allowMethod {
		if mlist == m {
			return true
		}
	}
	return ok
}
