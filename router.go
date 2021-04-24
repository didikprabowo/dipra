package dipra

import (
	"net/http"
)

var (
	// List method is allowed
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

func (r *route) getParams() *params {
	ps, _ := r.pool.Get().(*params)
	*ps = (*ps)[0:0]
	return ps
}
