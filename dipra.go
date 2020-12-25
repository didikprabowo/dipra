package dipra

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
)

var (
	welcome = `Welcome to Dipra mini framework golang for build service`
)

type (

	// HandlerFunc ro running HandlerFunc context
	HandlerFunc func(*Context) error

	// MiddlewareFunc to handle middleware
	MiddlewareFunc func(HandlerFunc) HandlerFunc

	// Engine core of dipra
	Engine struct {
		Prefix string
		// Route
		Route []Route
		// MiddlewareFunc func
		HandleMiddleware []MiddlewareFunc
		// Sync.Pool
		Pool sync.Pool
		// Node
		Node Node
	}

	// Route for handler routing
	Route struct {
		Path    string
		Method  string
		Handler HandlerFunc
	}

	// M map[string]interface{}
	M map[string]interface{}
)

const (
	// DefaultPort is 8080
	DefaultPort string = ":8080"
)

// Default Engine
func Default() *Engine {
	e := &Engine{
		Route: []Route{},
		Node:  Node{},
	}
	e.Pool.New = func() interface{} {
		return e.InitialContext(nil, nil)
	}
	return e
}

// InitialContext ...
func (e *Engine) InitialContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		ResponseWriter: w,
		Request:        r,
		Writen: ResponseWriter{
			Response:   w,
			statusCode: http.StatusOK,
		},
		Params: Param{},
	}
}

// AddToObjectEngine is used for set routing and middleware
func (e *Engine) AddToObjectEngine(path, method string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.HandleMiddleware = append(e.HandleMiddleware, middleware...)
	e.Route = append(e.Route, Route{Path: e.Prefix + path, Method: method, Handler: handler})
}

// Use is used for add handlefuncs
func (e *Engine) Use(middleware ...MiddlewareFunc) {
	e.HandleMiddleware = append(e.HandleMiddleware, middleware...)
}

// Group is used for grouped route
func (e *Engine) Group(group string, m ...MiddlewareFunc) *Engine {
	e.Prefix = group
	e.HandleMiddleware = append(e.HandleMiddleware, m...)
	return e
}

// GET is used HTTP Request with GET METHOD
func (e *Engine) GET(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddToObjectEngine(path, http.MethodGet, handler, middleware...)
}

// POST is used HTTP Request with POST METHOD
func (e *Engine) POST(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddToObjectEngine(path, http.MethodPost, handler, middleware...)
}

// PUT is used HTTP Request with PUT METHOD
func (e *Engine) PUT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddToObjectEngine(path, http.MethodPut, handler, middleware...)
}

// PATCH is used HTTP Request with PATCH METHOD
func (e *Engine) PATCH(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddToObjectEngine(path, http.MethodPatch, handler, middleware...)
}

// DELETE is used HTTP Request with DELETE METHOD
func (e *Engine) DELETE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddToObjectEngine(path, http.MethodDelete, handler, middleware...)
}

// OPTION is used HTTP Request with OPTION METHOD
func (e *Engine) OPTION(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddToObjectEngine(path, http.MethodOptions, handler, middleware...)
}

// TRACE is used HTTP Request with TRACE METHOD
func (e *Engine) TRACE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddToObjectEngine(path, http.MethodTrace, handler, middleware...)
}

// CONNECT is used HTTP Request with CONNECT METHOD
func (e *Engine) CONNECT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddToObjectEngine(path, http.MethodConnect, handler, middleware...)
}

// Static is used define http to get file type
func (e *Engine) Static(prefix, root string) {
	p := func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			cleanURL := path.Clean(c.Request.URL.String())
			name := path.Join(root, strings.ReplaceAll(cleanURL, prefix, ""))

			isExist, err := os.Stat(root)
			if err != nil || !isExist.IsDir() {
				return err
			}
			return c.File(name)
		}
	}
	e.Use(p)
}

// defaulterrorHttp is used set error default
func defaulterrorHTTP(w http.ResponseWriter, code int, err error) MiddlewareFunc {
	return func(c HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			return c.JSON(code, M{
				"error": err,
			})
		}
	}
}

// defaultErrorHandler is used default handler
func defaultErrorHandler(c HandlerFunc, werrx *WrapError) HandlerFunc {
	return func(c *Context) error {
		return c.JSON(werrx.Code, M{
			"error": werrx,
		})
	}
}

// HandlerError for handler context error
func (e *Engine) HandlerError(err error, c *Context) {

	var (
		r = M{}
	)

	eStr, ok := err.(*WrapError)
	if !ok {
		r["code"] = http.StatusInternalServerError
		r["message"] = err.Error()
	} else {
		r["code"] = eStr.Code
		r["message"] = eStr.Message.(string)
	}

	err = c.JSON(r["code"].(int), M{
		"error": r,
	})
}

// ServeHTTP is used run http server
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := e.Pool.Get().(*Context)
	ctx.Reset(w, r)
	rt, err := e.HandlerRoute(ctx)
	if err != nil {
		rt.Handler = defaultErrorHandler(rt.Handler, err)
	}

	if e.HandleMiddleware != nil {
		rt.Handler = e.WrapMiddleware(rt.Handler)
	}

	if err := rt.Handler(ctx); err != nil {
		e.HandlerError(err, ctx)
	}
}

// HandlerRoute is used running context http
func (e *Engine) HandlerRoute(c *Context) (r Route, werrx *WrapError) {

	sort.Slice(e.Route, func(i, j int) bool {
		return (e.Route[i].Path[1:len(e.Route[i].Path)] == c.URL.String()[1:len(c.URL.String())])
	})

	for _, rt := range e.Route {
		e.Node.SetNode(c, rt)
		url, err := e.Node.ReserverURI()
		if err != nil {
			werrx := &WrapError{
				Code:     http.StatusInternalServerError,
				Message:  err.Error(),
				Internal: http.StatusText(http.StatusInternalServerError),
			}
			return rt, werrx
		}
		uriCtx := c.URL.EscapedPath()[1:len(c.URL.EscapedPath())]
		if uriCtx == url {
			if c.Method != rt.Method {
				return rt, &WrapError{
					Code:    http.StatusMethodNotAllowed,
					Message: http.StatusText(http.StatusMethodNotAllowed),
				}
			}
			return rt, werrx
		}
	}

	return r, &WrapError{
		Code:    http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	}
}

// WrapMiddleware is used wrapping with returns HandlerFunc
func (e *Engine) WrapMiddleware(h HandlerFunc) HandlerFunc {
	for i := len(e.HandleMiddleware) - 1; i >= 0; i-- {
		if e.HandleMiddleware[i](h) != nil {
			h = e.HandleMiddleware[i](h)
		}
	}
	return h
}

// Run Server with HTTP
func (e *Engine) Run(addr string) (err error) {
	if addr == "" {
		addr = DefaultPort
	}

	fmt.Fprintf(os.Stdout, fmt.Sprintf("%v\nServer started with http %v[::]:%v %v\n", welcome, Green, addr, Reset))

	err = http.ListenAndServe(addr, e)
	if err != nil {
		return err
	}
	return err
}

// RunTLSaddr Server with HTTPS
func (e *Engine) RunTLSaddr(addr string, certFile, keyFile string) (err error) {
	if addr == "" {
		addr = DefaultPort
	}
	fmt.Fprintf(os.Stdout, fmt.Sprintf("%v\nServer started with https %v[::]:%v %v\n", welcome, Green, addr, Reset))
	err = http.ListenAndServeTLS(addr, certFile, keyFile, e)
	if err != nil {
		return err
	}

	return err
}
