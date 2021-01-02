package dipra

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	// List method is allowed
	allowMethod = []string{
		http.MethodGet,
		http.MethodDelete,
		http.MethodPost,
		http.MethodConnect,
		http.MethodPatch,
		http.MethodPut,
	}
)

// allowMethods for use check your define a method
func (e *Engine) allowMethods(m string) (ok bool) {
	for _, mlist := range allowMethod {
		if mlist == m {
			return true
		}
	}
	return ok
}

func (e *Engine) cleanPath(p string) string {
	return strings.ReplaceAll(p, "/", "")
}

// allowPath for use check your path
func (e *Engine) allowPath(p string) (ok bool) {
	p = e.cleanPath(p)

	if p == "" {
		return true
	}

	isAllow := regexp.MustCompile(`[a-zA-Z0-9]$`)
	ok = isAllow.MatchString(p)

	return ok
}

/*
	findRouter for used to checking existing path,
	method and handler
*/
func (e *Engine) findRouter(m, p string, h HandlerFunc) (exist bool) {

	if ok := e.allowMethods(m); !ok {
		log.Printf("Method %s not allow", m)
		return
	}

	if ok := e.allowPath(p); !ok {
		log.Printf("Path %s is invalid syntax", p)
		return
	}

	toLower := func(s string) string {
		return strings.ToLower(s)
	}

	for i := range e.getRoutes() {
		if toLower(e.Route[i].Method) == toLower(m) &&
			toLower(e.Route[i].Path) == toLower(p) {
			e.Route[i].Handler = h
			return true
		}
	}
	return exist
}

// getRoutes for get The routes you define
func (e *Engine) getRoutes() []Route {
	return e.Route
}
