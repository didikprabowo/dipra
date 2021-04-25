package dipra

import (
	"net/http"
)

type (
	Group interface {
		GET(path string, handler HandlerFunc, middleware ...MiddlewareFunc)
		POST(path string, handler HandlerFunc, middleware ...MiddlewareFunc)
		PUT(path string, handler HandlerFunc, middleware ...MiddlewareFunc)
		PATCH(path string, handler HandlerFunc, middleware ...MiddlewareFunc)
		DELETE(path string, handler HandlerFunc, middleware ...MiddlewareFunc)
		ANY(path string, handler HandlerFunc, middleware ...MiddlewareFunc)
		use(m ...MiddlewareFunc)
	}
	groupRoute struct {
		prefix     string
		middleware []MiddlewareFunc
		dipra      *Engine
	}
)

func NewGroupRoute(prefix string, dipra *Engine) Group {

	return &groupRoute{
		dipra:      dipra,
		prefix:     prefix,
		middleware: []MiddlewareFunc{},
	}
}

func (g *groupRoute) addGroupRoute(path string, method string, handler HandlerFunc, middleware ...MiddlewareFunc) {

	if path == "/" {
		path = ""
	}

	middleware = append(middleware, g.middleware...)

	fullpath := g.prefix + path
	g.dipra.AddRoute(fullpath, method, handler, middleware...)
	g.dipra.group = append(g.dipra.group, *g)
}

func (g *groupRoute) GET(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.addGroupRoute(path, http.MethodGet, handler, middleware...)
}

func (g *groupRoute) POST(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.addGroupRoute(path, http.MethodPost, handler, middleware...)
}

func (g *groupRoute) PUT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.addGroupRoute(path, http.MethodPut, handler, middleware...)
}

func (g *groupRoute) PATCH(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.addGroupRoute(path, http.MethodPatch, handler, middleware...)
}

func (g *groupRoute) DELETE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.addGroupRoute(path, http.MethodDelete, handler, middleware...)
}

func (g *groupRoute) ANY(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	for _, method := range allowMethod {
		g.addGroupRoute(path, method, handler, middleware...)
	}
}

func (g *groupRoute) use(m ...MiddlewareFunc) {
	g.middleware = append(g.middleware, m...)
}
